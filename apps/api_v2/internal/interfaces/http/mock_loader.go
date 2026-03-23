package http_interfaces

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

var mockBasePaths = []string{
	"internal/interfaces/http",
	"test/http/mockup-responses",
}

func mockHandler(relativePath string, statusCode int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		payload, err := loadMock(relativePath)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "mock_response_not_found",
				"message": err.Error(),
			})
		}

		if envelope, ok := payload.(map[string]any); ok {
			if body, hasBody := envelope["body"]; hasBody {
				if code, hasCode := envelope["statusCode"]; hasCode {
					if codeFloat, castOk := code.(float64); castOk {
						c.Status(int(codeFloat))
						return c.JSON(body)
					}
				}
				c.Status(statusCode)
				return c.JSON(body)
			}
		}

		c.Status(statusCode)
		return c.JSON(payload)
	}
}

func loadMock(relativePath string) (any, error) {
	var lastErr error

	for _, basePath := range mockBasePaths {
		path := filepath.Join(basePath, relativePath)
		bytes, err := os.ReadFile(path)
		if err != nil {
			lastErr = err
			continue
		}

		var out any
		if err := json.Unmarshal(bytes, &out); err != nil {
			return nil, fmt.Errorf("decode mock file %s: %w", path, err)
		}

		return out, nil
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("mock file not found")
	}

	return nil, fmt.Errorf("read mock file %s: %w", relativePath, lastErr)
}
