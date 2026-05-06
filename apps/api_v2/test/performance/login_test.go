package performance_tests

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	httputils "github.com/openlabun/CODER/apps/api_v2/test/http"
)

const loginConcurrency = 10

func TestLoginPerformance(t *testing.T) {
	process, app := httputils.StartConcurrentHTTPTestWithApp(t, "Login Performance")

	timestamp := time.Now().UnixNano()
	password := "Password123!"

	emails := make([]string, loginConcurrency)
	for i := 0; i < loginConcurrency; i++ {
		emails[i] = fmt.Sprintf("perf+%d+%d@test.com", timestamp, i)
	}
	accessTokens := make([]string, loginConcurrency)
	refreshTokens := make([]string, loginConcurrency)
	userIDs := make([]string, loginConcurrency)

	// [STEP 1] Register test users (critical)
	var regCounter int64 = -1
	process.StartConcurrentStep(1, "Registrar usuarios concurrentemente", func() (int, []byte, error) {
		i := atomic.AddInt64(&regCounter, 1)
		resp := httputils.PostAuthRegister(t, app, emails[i], fmt.Sprintf("perf-%d", i), password)
		return resp.StatusCode, resp.Body, nil
	})
	process.RunStep(loginConcurrency)

	// [STEP 2] Login test users (critical)
	var loginCounter int64 = -1
	process.StartConcurrentStep(2, "Login concurrente", func() (int, []byte, error) {
		i := atomic.AddInt64(&loginCounter, 1)
		resp := httputils.PostAuthLogin(t, app, emails[i], password)
		if resp.StatusCode == 200 && resp.JSON != nil {
			if tokenObj, ok := resp.JSON["token"].(map[string]any); ok {
				accessTokens[i], _ = tokenObj["access_token"].(string)
				refreshTokens[i], _ = tokenObj["refresh_token"].(string)
			}
			if userObj, ok := resp.JSON["user_data"].(map[string]any); ok {
				userIDs[i], _ = userObj["id"].(string)
			}
		}
		return resp.StatusCode, resp.Body, nil
	})
	process.RunStep(loginConcurrency)

	// [STEP 3] Refresh token test users (critical)
	var refreshCounter int64 = -1
	process.StartConcurrentStep(3, "Refresh token concurrente", func() (int, []byte, error) {
		i := atomic.AddInt64(&refreshCounter, 1)
		if refreshTokens[i] == "" {
			return 0, nil, fmt.Errorf("no refresh token for user %d", i)
		}
		resp := httputils.PostAuthRefreshToken(t, app, refreshTokens[i])
		return resp.StatusCode, resp.Body, nil
	})
	process.RunStep(loginConcurrency)

	// [STEP 4] Get user info test users (critical)
	var meCounter int64 = -1
	process.StartConcurrentStep(4, "Get /auth/me concurrente", func() (int, []byte, error) {
		i := atomic.AddInt64(&meCounter, 1)
		if accessTokens[i] == "" {
			return 0, nil, fmt.Errorf("no access token for user %d", i)
		}
		access := &httputils.HTTPAccess{
			UserID:      userIDs[i],
			Email:       emails[i],
			AccessToken: accessTokens[i],
		}
		resp := httputils.GetAuthMe(t, app, access, "")
		return resp.StatusCode, resp.Body, nil
	})
	process.RunStep(loginConcurrency)

	process.End()
}
