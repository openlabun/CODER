package utils

import (
	"bytes"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"

	container "github.com/openlabun/CODER/apps/api_v2/internal/application/container"
	http_interfaces "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http"
)

type requestMock struct {
	Method     string         `json:"method"`
	Path       string         `json:"path"`
	PathParams map[string]any `json:"pathParams"`
	Query      map[string]any `json:"query"`
	Body       any            `json:"body"`
}

type responseMock struct {
	StatusCode int `json:"statusCode"`
	Body       any `json:"body"`
}

func ReadFileContent(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func ReadInputFileContent(module, endpoint string) (string, error) {
	path := filepath.Join(apiV2RootDir(), "internal", "interfaces", "http", module, endpoint, "mockup", "input.json")
	return ReadFileContent(path)
}

func ReadOutputFileContent(module, endpoint string) (string, error) {
	path := filepath.Join(apiV2RootDir(), "internal", "interfaces", "http", module, endpoint, "mockup", "output.json")
	return ReadFileContent(path)
}

func InitApp() (*fiber.App, error) {
	moduleRoot := apiV2RootDir()
	_ = loadEnvFile(filepath.Join(moduleRoot, ".env.dev"))
	_ = loadEnvFile(filepath.Join(moduleRoot, "..", "..", ".env.dev"))

	if strings.TrimSpace(os.Getenv("ROBLE_PROJECT")) == "" {
		_ = os.Setenv("ROBLE_PROJECT", "test-project")
	}
	if strings.TrimSpace(os.Getenv("ROBLE_BASE_URL")) == "" {
		_ = os.Setenv("ROBLE_BASE_URL", "https://roble.invalid")
	}

	appContainer, err := container.BuildApplicationContainer()
	if err != nil {
		return nil, fmt.Errorf("initialize application container: %w", err)
	}

	app := fiber.New()
	http_interfaces.RegisterRoutes(app, appContainer)

	return app, nil
}

func loadEnvFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, "\"")
		value = strings.Trim(value, "'")

		if key == "" {
			continue
		}

		if strings.TrimSpace(os.Getenv(key)) == "" {
			_ = os.Setenv(key, value)
		}
	}

	return scanner.Err()
}

func getRequest(method, url string, headers map[string]string, params []byte) (*fiber.Ctx, error) {
	return runRequest(method, url, headers, params)
}

func postRequest(url string, headers map[string]string, body []byte) (*fiber.Ctx, error) {
	return runRequest(http.MethodPost, url, headers, body)
}

func putRequest(url string, headers map[string]string, body []byte) (*fiber.Ctx, error) {
	return runRequest(http.MethodPut, url, headers, body)
}

func deleteRequest(url string, headers map[string]string, params []byte) (*fiber.Ctx, error) {
	return runRequest(http.MethodDelete, url, headers, params)
}

func CompareResponse(expected, actual string) bool {
	var expectedJSON any
	if err := json.Unmarshal([]byte(expected), &expectedJSON); err != nil {
		return strings.TrimSpace(expected) == strings.TrimSpace(actual)
	}

	var actualJSON any
	if err := json.Unmarshal([]byte(actual), &actualJSON); err != nil {
		return false
	}

	left, err := json.Marshal(expectedJSON)
	if err != nil {
		return false
	}
	right, err := json.Marshal(actualJSON)
	if err != nil {
		return false
	}

	return bytes.Equal(left, right)
}

func RunEndpointMockupTest(t interface {
	Helper()
	Logf(format string, args ...any)
	Fatalf(format string, args ...any)
	Cleanup(func())
}, module, endpoint string) {
	t.Helper()

	inputContent, err := ReadInputFileContent(module, endpoint)
	if err != nil {
		t.Fatalf("read input.json failed for %s/%s: %v", module, endpoint, err)
	}

	outputContent, err := ReadOutputFileContent(module, endpoint)
	if err != nil {
		t.Fatalf("read output.json failed for %s/%s: %v", module, endpoint, err)
	}

	var input requestMock
	if err := json.Unmarshal([]byte(inputContent), &input); err != nil {
		t.Fatalf("decode input.json failed for %s/%s: %v", module, endpoint, err)
	}

	var expected responseMock
	if err := json.Unmarshal([]byte(outputContent), &expected); err != nil {
		t.Fatalf("decode output.json failed for %s/%s: %v", module, endpoint, err)
	}

	requestPath := buildRequestPath(input.Path, input.PathParams, input.Query)
	bodyBytes := []byte(nil)
	if input.Body != nil {
		bodyBytes, err = json.Marshal(input.Body)
		if err != nil {
			t.Fatalf("encode request body failed for %s/%s: %v", module, endpoint, err)
		}
	}

	t.Logf("CHECK REQUEST module=%s endpoint=%s method=%s path=%s body=%s", module, endpoint, strings.ToUpper(input.Method), requestPath, string(bodyBytes))
	t.Logf("CHECK EXPECTED module=%s endpoint=%s status=%d response=%s", module, endpoint, expected.StatusCode, outputContent)

	ctx, err := getRequest(strings.ToUpper(input.Method), requestPath, map[string]string{"Content-Type": "application/json"}, bodyBytes)
	if err != nil {
		t.Fatalf("http request failed for %s/%s: %v", module, endpoint, err)
	}
	t.Cleanup(func() {
		ctx.App().ReleaseCtx(ctx)
	})

	actualStatus := ctx.Response().StatusCode()
	actualBody := string(ctx.Response().Body())

	t.Logf("CHECK ACTUAL module=%s endpoint=%s status=%d response=%s", module, endpoint, actualStatus, actualBody)

	if actualStatus != expected.StatusCode {
		t.Fatalf("unexpected status for %s/%s: got=%d expected=%d", module, endpoint, actualStatus, expected.StatusCode)
	}

	expectedBodyBytes, err := json.Marshal(expected.Body)
	if err != nil {
		t.Fatalf("encode expected body failed for %s/%s: %v", module, endpoint, err)
	}

	if !CompareResponse(string(expectedBodyBytes), actualBody) {
		t.Fatalf("unexpected response body for %s/%s: got=%s expected=%s", module, endpoint, actualBody, string(expectedBodyBytes))
	}
}

func runRequest(method, requestURL string, headers map[string]string, body []byte) (*fiber.Ctx, error) {
	app, err := InitApp()
	if err != nil {
		return nil, err
	}

	baseDir := apiV2RootDir()
	originalWD, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	if err := os.Chdir(baseDir); err != nil {
		return nil, err
	}
	defer func() {
		_ = os.Chdir(originalWD)
	}()

	req := httptest.NewRequest(method, requestURL, bytes.NewReader(body))
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := app.Test(req, int((10 * time.Second).Milliseconds()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Response().SetStatusCode(resp.StatusCode)
	ctx.Response().SetBody(respBody)

	return ctx, nil
}

func buildRequestPath(pathTemplate string, pathParams map[string]any, query map[string]any) string {
	path := pathTemplate
	for key, value := range pathParams {
		token := "{" + key + "}"
		path = strings.ReplaceAll(path, token, url.PathEscape(fmt.Sprint(value)))
	}

	if len(query) == 0 {
		return path
	}

	values := url.Values{}
	for key, value := range query {
		values.Set(key, fmt.Sprint(value))
	}

	encoded := values.Encode()
	if encoded == "" {
		return path
	}

	if strings.Contains(path, "?") {
		return path + "&" + encoded
	}

	return path + "?" + encoded
}

func apiV2RootDir() string {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		return "."
	}
	return filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", ".."))
}

func DoJSONRequest(app *fiber.App, method, path string, body any, headers map[string]string) (int, []byte, error) {
	var requestBody []byte
	var err error

	if body != nil {
		requestBody, err = json.Marshal(body)
		if err != nil {
			return 0, nil, fmt.Errorf("marshal request body: %w", err)
		}
	}

	req := httptest.NewRequest(method, path, bytes.NewReader(requestBody))
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := app.Test(req, 10000)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, fmt.Errorf("read response body: %w", err)
	}

	return resp.StatusCode, respBody, nil
}
