package utils

import (
	"bufio"
	"bytes"
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
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"

	container "github.com/openlabun/CODER/apps/api_v2/internal/application/container"
	http_interfaces "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http"

	dtos "github.com/openlabun/CODER/apps/api_v2/internal/application/dtos/user"
	Entities "github.com/openlabun/CODER/apps/api_v2/internal/domain/entities/user"
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

func EnsureHTTPAuthUserAccess(t *testing.T, app *fiber.App, email, password, name string) *dtos.UserAccess {
	t.Helper()

	t.Logf("[STEP] Attempt HTTP login for %s", email)
	loginBody := map[string]string{
		"email":    email,
		"password": password,
	}

	status, body, err := PostAuthLogin(nil, loginBody)
	if err != nil {
		t.Fatalf("login request failed for %s: %v", email, err)
	}

	if status == http.StatusOK {
		access := decodeUserAccess(t, body, "login")
		validateUserAccess(t, access, email)
		return access
	}

	t.Logf("login failed for %s with status=%d, trying register", email, status)
	registerBody := map[string]string{
		"email":    email,
		"name":     name,
		"password": password,
	}

	status, body, err = PostAuthRegister(nil, registerBody)
	if err != nil {
		t.Fatalf("register request failed for %s: %v", email, err)
	}
	if status != http.StatusCreated {
		t.Fatalf("expected register status=%d for %s, got=%d body=%s", http.StatusCreated, email, status, string(body))
	}

	access := decodeUserAccess(t, body, "register")
	validateUserAccess(t, access, email)
	return access
}

func GetUserDataHTTP(t *testing.T, app *fiber.App, access *dtos.UserAccess) *Entities.User {
	t.Helper()

	headers := map[string]string{
		"Authorization": "Bearer " + access.Token.AccessToken,
		"X-User-Email":  access.UserData.Email,
	}

	status, body, err := GetAuthMe(headers)
	if err != nil {
		t.Fatalf("/auth/me request failed: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("expected /auth/me status=%d, got=%d body=%s", http.StatusOK, status, string(body))
	}

	var user Entities.User
	if err := json.Unmarshal(body, &user); err != nil {
		t.Fatalf("decode /auth/me response failed: %v body=%s", err, string(body))
	}

	if user.ID == "" {
		t.Fatal("expected /auth/me response with user ID")
	}

	return &user
}

func decodeUserAccess(t *testing.T, raw []byte, source string) *dtos.UserAccess {
	t.Helper()

	var access dtos.UserAccess
	if err := json.Unmarshal(raw, &access); err != nil {
		t.Fatalf("decode %s response failed: %v body=%s", source, err, string(raw))
	}

	return &access
}

func validateUserAccess(t *testing.T, access *dtos.UserAccess, expectedEmail string) {
	t.Helper()

	if access == nil || access.UserData == nil || access.Token == nil {
		t.Fatalf("expected valid access payload for %s", expectedEmail)
	}
	if access.UserData.ID == "" {
		t.Fatalf("expected user ID in access payload for %s", expectedEmail)
	}
	if access.Token.AccessToken == "" {
		t.Fatalf("expected access token in payload for %s", expectedEmail)
	}
	if access.UserData.Email != expectedEmail {
		t.Fatalf("expected email=%s in payload, got=%s", expectedEmail, access.UserData.Email)
	}
}
