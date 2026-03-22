package roble_infrastructure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	defaultRobleBaseURL = "https://roble-api.openlab.uninorte.edu.co"
)

type RobleClient struct {
	httpClient *http.Client
	baseURL    string
	project    string
}

func NewRobleClient(httpClient *http.Client) (*RobleClient, error) {
	project := strings.TrimSpace(os.Getenv("ROBLE_PROJECT"))
	if project == "" {
		return nil, fmt.Errorf("ROBLE_PROJECT is required")
	}

	baseURL := strings.TrimSpace(os.Getenv("ROBLE_BASE_URL"))
	if baseURL == "" {
		baseURL = defaultRobleBaseURL
	}

	return NewRobleClientWithConfig(baseURL, project, httpClient)
}

func NewRobleClientWithConfig(baseURL, project string, httpClient *http.Client) (*RobleClient, error) {
	if strings.TrimSpace(project) == "" {
		return nil, fmt.Errorf("roble project is required")
	}

	if strings.TrimSpace(baseURL) == "" {
		return nil, fmt.Errorf("roble base URL is required")
	}

	if httpClient == nil {
		httpClient = &http.Client{Timeout: 15 * time.Second}
	}

	return &RobleClient{
		httpClient: httpClient,
		baseURL:    strings.TrimRight(baseURL, "/"),
		project:    project,
	}, nil
}

func (c *RobleClient) Login(email, password string) (*RobleLoginResponse, error) {
	body := map[string]string{
		"email":    email,
		"password": password,
	}

	var out RobleLoginResponse
	if err := c.doJSON(http.MethodPost, c.authURL("login"), "", body, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

func (c *RobleClient) Signup(email, password, name string) (*RobleMessageResponse, error) {
	body := map[string]string{
		"email":    email,
		"password": password,
		"name":     name,
	}

	var out RobleMessageResponse
	if err := c.doJSONExpectedStatus(http.MethodPost, c.authURL("signup"), "", body, &out, http.StatusCreated); err != nil {
		return nil, err
	}

	return &out, nil
}

func (c *RobleClient) SignupDirect(email, password, name string) (*RobleMessageResponse, error) {
	body := map[string]string{
		"email":    email,
		"password": password,
		"name":     name,
	}

	var out RobleMessageResponse
	if err := c.doJSONExpectedStatus(http.MethodPost, c.authURL("signup-direct"), "", body, &out, http.StatusCreated); err != nil {
		return nil, err
	}

	return &out, nil
}

func (c *RobleClient) RefreshToken(refreshToken string) (*RobleRefreshTokenResponse, error) {
	body := map[string]string{
		"refreshToken": refreshToken,
	}

	var out RobleRefreshTokenResponse
	if err := c.doJSON(http.MethodPost, c.authURL("refresh-token"), "", body, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

func (c *RobleClient) VerifyEmail(email, code string) error {
	body := map[string]string{
		"email": email,
		"code":  code,
	}

	return c.doJSON(http.MethodPost, c.authURL("verify-email"), "", body, nil)
}

func (c *RobleClient) ForgotPassword(email string) error {
	body := map[string]string{
		"email": email,
	}

	return c.doJSON(http.MethodPost, c.authURL("forgot-password"), "", body, nil)
}

func (c *RobleClient) ResetPassword(token, newPassword string) error {
	body := map[string]string{
		"token":       token,
		"newPassword": newPassword,
	}

	return c.doJSON(http.MethodPost, c.authURL("reset-password"), "", body, nil)
}

func (c *RobleClient) Logout(accessToken string) error {
	return c.doJSON(http.MethodPost, c.authURL("logout"), accessToken, nil, nil)
}

func (c *RobleClient) VerifyToken(accessToken string) (*RobleVerifyTokenResponse, error) {
	var out RobleVerifyTokenResponse
	if err := c.doJSON(http.MethodGet, c.authURL("verify-token"), accessToken, nil, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

func (c *RobleClient) Insert(tableName string, records []map[string]any, accessToken string) (map[string]any, error) {
	var out map[string]any
	err := c.doJSON(
		http.MethodPost,
		c.databaseURL("insert"),
		accessToken,
		RobleInsertRequest{TableName: tableName, Records: records},
		&out,
	)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (c *RobleClient) Read(tableName string, conditions map[string]string, accessToken string) (map[string]any, error) {
	query := map[string]string{"tableName": tableName}
	for k, v := range conditions {
		query[k] = v
	}

	var raw any
	err := c.doJSONWithQuery(http.MethodGet, c.databaseURL("read"), accessToken, query, nil, &raw)
	if err != nil {
		return nil, err
	}

	switch out := raw.(type) {
	case map[string]any:
		return out, nil
	case []any:
		// Normalize array responses so repository layer can consume a stable shape.
		return map[string]any{"data": out}, nil
	default:
		return nil, fmt.Errorf("unexpected read response type %T", raw)
	}
}

func (c *RobleClient) Update(tableName, idColumn, idValue string, updates map[string]any, accessToken string) (map[string]any, error) {
	var out map[string]any
	err := c.doJSON(
		http.MethodPut,
		c.databaseURL("update"),
		accessToken,
		RobleUpdateRequest{TableName: tableName, IDColumn: idColumn, IDValue: idValue, Updates: updates},
		&out,
	)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (c *RobleClient) Delete(tableName, idColumn, idValue, accessToken string) (map[string]any, error) {
	var out map[string]any
	err := c.doJSON(
		http.MethodDelete,
		c.databaseURL("delete"),
		accessToken,
		RobleDeleteRequest{TableName: tableName, IDColumn: idColumn, IDValue: idValue},
		&out,
	)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (c *RobleClient) authURL(path string) string {
	return fmt.Sprintf("%s/auth/%s/%s", c.baseURL, c.project, strings.TrimLeft(path, "/"))
}

func (c *RobleClient) databaseURL(path string) string {
	return fmt.Sprintf("%s/database/%s/%s", c.baseURL, c.project, strings.TrimLeft(path, "/"))
}

func (c *RobleClient) doJSONWithQuery(method, rawURL, accessToken string, query map[string]string, requestBody any, output any) error {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("parse url: %w", err)
	}

	q := parsedURL.Query()
	for k, v := range query {
		q.Set(k, v)
	}
	parsedURL.RawQuery = q.Encode()

	return c.doJSON(method, parsedURL.String(), accessToken, requestBody, output)
}

func (c *RobleClient) doJSON(method, url, accessToken string, requestBody any, output any) error {
	return c.doJSONExpectedStatus(method, url, accessToken, requestBody, output)
}

func (c *RobleClient) doJSONExpectedStatus(method, url, accessToken string, requestBody any, output any, expectedStatus ...int) error {
	var reader io.Reader
	if requestBody != nil {
		payload, err := json.Marshal(requestBody)
		if err != nil {
			return fmt.Errorf("marshal request body: %w", err)
		}
		reader = bytes.NewBuffer(payload)
	}

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if strings.TrimSpace(accessToken) != "" {
		req.Header.Set("Authorization", "Bearer "+strings.TrimSpace(accessToken))
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request to roble failed: %w", err)
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	statusValid := false
	if len(expectedStatus) == 0 {
		statusValid = res.StatusCode >= 200 && res.StatusCode < 300
	} else {
		for _, status := range expectedStatus {
			if res.StatusCode == status {
				statusValid = true
				break
			}
		}
	}


	if !statusValid {
		bodyText := strings.TrimSpace(string(bodyBytes))
		if len(bodyText) > 300 {
			bodyText = bodyText[:300]
		}
		if bodyText == "" {
			bodyText = "empty response body"
		}
		return fmt.Errorf("roble request failed (%d): %s", res.StatusCode, bodyText)
	}

	if res.StatusCode == http.StatusCreated {
		trimmedBody := bytes.TrimSpace(bodyBytes)
		if len(trimmedBody) > 0 {
			var parsedBody any
			if err := json.Unmarshal(trimmedBody, &parsedBody); err == nil {
				if bodyMap, ok := parsedBody.(map[string]any); ok {
					if skippedRaw, exists := bodyMap["skipped"]; exists {
						skippedItems, isArray := skippedRaw.([]any)
						if !isArray {
							return fmt.Errorf("roble request failed (201): skipped must be an array")
						}
						if len(skippedItems) > 0 {
							return fmt.Errorf("roble request failed (201): database returned skipped items: %v", skippedItems)
						}
					}
				}
			}
		}
	}

	
	if output == nil || len(bytes.TrimSpace(bodyBytes)) == 0 {
		return nil
	}

	if err := json.Unmarshal(bodyBytes, output); err != nil {
		return fmt.Errorf("decode response body: %w", err)
	}

	

	return nil
}
