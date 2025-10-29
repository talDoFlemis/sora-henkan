package http

import (
	"net/http"
	"strings"

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
