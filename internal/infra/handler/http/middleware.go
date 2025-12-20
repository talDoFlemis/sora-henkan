package http

import (
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func CloudFlareExtractClientIPfunc(req *http.Request) string {
	cf := req.Header["CF-Connecting-IP"]
	if len(cf) > 0 {
		return cf[0]
	}

	return echo.ExtractIPFromXFFHeader()(req)
}

func SSETimeoutSkipper(c echo.Context) bool {
	// Skip timeout for SSE endpoints
	path := c.Path()

	if c.Request().Header.Get("Accept") == "text/event-stream" || strings.Contains(path, "/sse") {
		return true
	}

	return false
}

func ValidationErrorMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)

			// Don't handle if response already started
			if c.Response().Committed {
				return nil
			}

			if validationErr, ok := err.(validator.ValidationErrors); ok {
				HandleValidationError(c, validationErr)
				return nil
			}

			he, ok := err.(*echo.HTTPError)
			if !ok {
				return c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"message": "Internal server error",
					"error":   err.Error(),
				})
			}

			// Handle Echo HTTP errors
			code := he.Code
			message := he.Message

			// If message is a string, use it directly
			if msg, ok := message.(string); ok {
				return c.JSON(code, map[string]interface{}{
					"message": msg,
				})
			}

			return c.JSON(code, message)
		}
	}
}
