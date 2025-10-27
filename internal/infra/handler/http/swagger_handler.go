package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/taldoflemis/sora-henkan/docs" // swagger docs
)

type SwaggerHandler struct{}

func NewSwaggerHandler() *SwaggerHandler {
	return &SwaggerHandler{}
}

func (h *SwaggerHandler) RegisterRoute(e *echo.Echo) {
	e.GET("/swagger/*", h.ServeSwaggerUI)
	e.GET("/docs/swagger.json", h.ServeSwaggerJSON)
}

// ServeSwaggerUI serves the Scalar UI for Swagger documentation
func (h *SwaggerHandler) ServeSwaggerUI(c echo.Context) error {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>API Documentation</title>
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1">
</head>
<body>
    <script id="api-reference" data-url="/docs/swagger.json"></script>
    <script>
      var configuration = {
        theme: 'purple',
      }

      document.getElementById('api-reference').dataset.configuration =
        JSON.stringify(configuration)
    </script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
</body>
</html>`
	return c.HTML(http.StatusOK, html)
}

// ServeSwaggerJSON serves the swagger.json file
func (h *SwaggerHandler) ServeSwaggerJSON(c echo.Context) error {
	return c.File("./docs/swagger.json")
}
