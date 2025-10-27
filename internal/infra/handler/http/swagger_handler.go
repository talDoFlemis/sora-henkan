package http

import (
	"net/http"

	scalargo "github.com/bdpiprava/scalar-go"
	"github.com/labstack/echo/v4"
)

type SwaggerHandler struct{}

func NewSwaggerHandler() *SwaggerHandler {
	return &SwaggerHandler{}
}

func (h *SwaggerHandler) RegisterRoute(e *echo.Echo) {
	e.GET("/swagger/*", h.ServeSwaggerUI)
	e.GET("/docs/openapi.yaml", h.ServeOpenAPISpec)
	// Keep backward compatibility
	e.GET("/docs/swagger.json", h.ServeSwaggerJSON)
}

// ServeSwaggerUI serves the Scalar UI for API documentation
func (h *SwaggerHandler) ServeSwaggerUI(c echo.Context) error {
	htmlContent, err := scalargo.NewV2(
		scalargo.WithSpecURL("/docs/openapi.yaml"),
		scalargo.WithDarkMode(),
		scalargo.WithMetaDataOpts(
			scalargo.WithTitle("Sora Henkan API Documentation"),
		),
	)
	
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to generate documentation")
	}
	
	return c.HTML(http.StatusOK, htmlContent)
}

// ServeOpenAPISpec serves the openapi.yaml file
func (h *SwaggerHandler) ServeOpenAPISpec(c echo.Context) error {
	return c.File("./docs/openapi.yaml")
}

// ServeSwaggerJSON serves the swagger.json file for backward compatibility
func (h *SwaggerHandler) ServeSwaggerJSON(c echo.Context) error {
	return c.File("./docs/swagger.json")
}
