package http

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/taldoflemis/sora-henkan/internal/core/application"
	"github.com/taldoflemis/sora-henkan/internal/core/domain/images"
	"github.com/taldoflemis/sora-henkan/internal/core/ports"
)

type ImageHandler struct {
	imageUseCase *application.ImageUseCase
}

func NewImageHandler(imageUseCase *application.ImageUseCase) *ImageHandler {
	return &ImageHandler{
		imageUseCase: imageUseCase,
	}
}

func (h *ImageHandler) RegisterRoute(g *echo.Group) {
	imageHandlerGroup := g.Group("v1/images")

	imageHandlerGroup.GET("/", h.ListImages)
	imageHandlerGroup.GET("/sse", h.GetAllImagesRealtimeUpdates)
	imageHandlerGroup.GET("/:id/sse", h.GetImageRealtimeUpdate)
	imageHandlerGroup.GET("/:id", h.GetImage)
	imageHandlerGroup.POST("/", h.CreateImage)
	imageHandlerGroup.PUT("/", h.UpdateImage)
	imageHandlerGroup.DELETE("/:id", h.DeleteImage)
}

func (h *ImageHandler) ListImages(c echo.Context) error {
	ctx := c.Request().Context()

	req := images.ListImagesRequest{}

	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	resp, err := h.imageUseCase.ListImages(ctx, &req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *ImageHandler) GetAllImagesRealtimeUpdates(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{})
}

func (h *ImageHandler) GetImageRealtimeUpdate(c echo.Context) error {
	ctx := context.Background()
	flusher, ok := c.Response().Writer.(http.Flusher)
	if !ok {
		slog.ErrorContext(ctx, "streaming unsupported by response writer")
		return echo.NewHTTPError(http.StatusInternalServerError, "Streaming unsupported")
	}

	id := c.Param("id")

	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")

	notify := c.Request().Context().Done()

	imageUpdates, closeCallback, err := h.imageUseCase.GetImageRealtimeUpdate(ctx, id)
	if err != nil {
		if errors.Is(err, ports.ErrImageNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "image not found"})
		}

		return err
	}
	defer closeCallback()

	for {
		select {
		case <-notify:
			slog.InfoContext(ctx, "client closed connection")
			return nil
		case update := <-imageUpdates:
			data, err := json.Marshal(update)
			if err != nil {
				slog.ErrorContext(ctx, "marshal order for SSE", slog.String("error", err.Error()))
				continue
			}
			_, err = c.Response().Writer.Write([]byte("data: " + string(data) + "\n\n"))
			if err != nil {
				slog.ErrorContext(ctx, "write SSE", slog.String("error", err.Error()))
				return err
			}
			flusher.Flush()
		}
	}
}

func (h *ImageHandler) GetImage(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")

	image, err := h.imageUseCase.GetImage(ctx, id)
	if err != nil {
		if errors.Is(err, ports.ErrImageNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "image not found"})
		}

		return err
	}

	return c.JSON(http.StatusOK, image)
}

func (h *ImageHandler) CreateImage(c echo.Context) error {
	ctx := c.Request().Context()
	req := images.CreateImageRequest{}

	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	resp, err := h.imageUseCase.CreateImageRequest(ctx, &req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, resp)
}

func (h *ImageHandler) UpdateImage(c echo.Context) error {
	ctx := c.Request().Context()
	req := images.UpdateImageRequest{}

	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	err = h.imageUseCase.UpdateImage(ctx, &req)
	if err != nil {
		if errors.Is(err, ports.ErrImageNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "image not found"})
		}

		return err
	}

	return c.JSON(http.StatusOK, map[string]string{})
}

func (h *ImageHandler) DeleteImage(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")

	err := h.imageUseCase.DeleteImage(ctx, id)
	if err != nil {
		if errors.Is(err, ports.ErrImageNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "image not found"})
		}

		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "image deleted successfully"})
}
