package http

import (
	"net/http"

	healthgo "github.com/hellofresh/health-go/v5"
	"github.com/labstack/echo/v4"
)

type HealthHandler struct {
	health *healthgo.Health
}

func NewHealthHandler(health *healthgo.Health) *HealthHandler {
	return &HealthHandler{
		health: health,
	}
}

func (h *HealthHandler) RegisterRoute(g *echo.Group) {
	g.GET("healthz", h.Handle)
}

func (h *HealthHandler) Handle(c echo.Context) error {
	check := h.health.Measure(c.Request().Context())

	statusCode := http.StatusOK
	if check.Status != healthgo.StatusOK {
		statusCode = http.StatusServiceUnavailable
	}

	return c.JSON(statusCode, check)
}
