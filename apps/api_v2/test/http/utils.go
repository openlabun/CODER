package http_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/gofiber/fiber/v2"
	http_interfaces "github.com/openlabun/CODER/apps/api_v2/internal/interfaces/http"
	test "github.com/openlabun/CODER/apps/api_v2/test"
)

type HTTPAccess struct {
	UserID      string
	Email       string
	AccessToken string
}

type HTTPResponse struct {
	StatusCode int
	Body       []byte
	JSON       map[string]any
}

func StartHTTPTestWithApp(t *testing.T, name string) (test.TestProcess, *fiber.App) {
	process := test.StartTestWithApp(t, name)
	app := fiber.New()
	http_interfaces.RegisterRoutes(app, process.Application)
	return process, app
}

func EnsureAuthUserAccess(t *testing.T, app *fiber.App, email, password, name string) *HTTPAccess {
	t.Helper()

	loginResp := PostAuthLogin(t, app, email, password)
	if loginResp.StatusCode == 200 {
		return ParseAccessResponse(t, loginResp, email)
	}

	registerResp := PostAuthRegister(t, app, email, name, password)
	if registerResp.StatusCode != 201 {
		t.Fatalf("register failed for %s: status=%d body=%s", email, registerResp.StatusCode, string(registerResp.Body))
	}

	return ParseAccessResponse(t, registerResp, email)
}

func ParseAccessResponse(t *testing.T, resp *HTTPResponse, expectedEmail string) *HTTPAccess {
	t.Helper()
	body := MustJSONMap(t, resp)

	tokenRaw, ok := body["token"].(map[string]any)
	if !ok {
		t.Fatalf("expected token object in response: %s", string(resp.Body))
	}

	userRaw, ok := body["user_data"].(map[string]any)
	if !ok {
		t.Fatalf("expected user_data object in response: %s", string(resp.Body))
	}

	userID := StringField(userRaw, "id")
	email := StringField(userRaw, "email")
	accessToken := StringField(tokenRaw, "access_token")

	if userID == "" || email == "" || accessToken == "" {
		t.Fatalf("expected non-empty auth payload fields: %s", string(resp.Body))
	}
	if expectedEmail != "" && email != expectedEmail {
		t.Fatalf("expected email=%s, got=%s", expectedEmail, email)
	}

	return &HTTPAccess{UserID: userID, Email: email, AccessToken: accessToken}
}

func MustJSONMap(t *testing.T, resp *HTTPResponse) map[string]any {
	t.Helper()
	if resp.JSON == nil {
		t.Fatalf("expected JSON map response, got empty body")
	}
	return resp.JSON
}

func RequireStatus(t *testing.T, resp *HTTPResponse, expected int, section string) {
	t.Helper()
	if resp.StatusCode != expected {
		t.Fatalf("%s: expected status=%d, got=%d body=%s", section, expected, resp.StatusCode, string(resp.Body))
	}
}

func StringField(obj map[string]any, key string) string {
	value, _ := obj[key]
	str, _ := value.(string)
	return str
}

func doJSONRequest(t *testing.T, app *fiber.App, method, path string, body any, headers map[string]string) *HTTPResponse {
	t.Helper()

	var payload []byte
	var err error
	if body != nil {
		payload, err = json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal request body for %s %s: %v", method, path, err)
		}
	}

	req := httptest.NewRequest(method, path, bytes.NewReader(payload))
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	timeoutMS := 30000
	if raw := os.Getenv("HTTP_TEST_TIMEOUT"); raw != "" {
		parsed, parseErr := strconv.Atoi(raw)
		if parseErr == nil && parsed > 0 {
			timeoutMS = parsed
		}
	}

	res, err := app.Test(req, timeoutMS)
	if err != nil {
		t.Fatalf("http request failed for %s %s: %v", method, path, err)
	}

	bodyBytes := make([]byte, res.ContentLength)
	if res.ContentLength <= 0 {
		bodyBytes = nil
	}
	if res.Body != nil {
		defer res.Body.Close()
		bodyBytes, err = io.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("read response body for %s %s: %v", method, path, err)
		}
	}

	response := &HTTPResponse{StatusCode: res.StatusCode, Body: bodyBytes}
	if len(bodyBytes) > 0 {
		var jsonBody map[string]any
		if json.Unmarshal(bodyBytes, &jsonBody) == nil {
			response.JSON = jsonBody
		}
	}

	return response
}

func authHeaders(access *HTTPAccess) map[string]string {
	if access == nil {
		return map[string]string{}
	}
	return map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", access.AccessToken),
		"X-User-Email":  access.Email,
	}
}
