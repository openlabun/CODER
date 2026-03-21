package http_interfaces

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

const mockBasePath = "test/http/mockup-responses"

func mockHandler(relativePath string, statusCode int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		payload, err := loadMock(relativePath)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "mock_response_not_found",
				"message": err.Error(),
			})
		}

		c.Status(statusCode)
		return c.JSON(payload)
	}
}

func loadMock(relativePath string) (any, error) {
	path := filepath.Join(mockBasePath, relativePath)
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read mock file %s: %w", path, err)
	}

	var out any
	if err := json.Unmarshal(bytes, &out); err != nil {
		return nil, fmt.Errorf("decode mock file %s: %w", path, err)
	}

	return out, nil
}
