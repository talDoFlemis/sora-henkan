package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func CloudFlareExtractClientIPfunc(req *http.Request) string {
	cf := req.Header["CF-Connecting-IP"]
	if len(cf) > 0 {
		return cf[0]
	}

	return echo.ExtractIPFromXFFHeader()(req)
}
