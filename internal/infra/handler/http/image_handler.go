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

// ListImages godoc
//
//	@Summary		List images
//	@Description	Get a paginated list of images
//	@Tags			images
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int	false	"Page number"		default(1)	minimum(1)
//	@Param			limit	query		int	false	"Items per page"	default(10)	minimum(1)	maximum(100)
//	@Success		200		{object}	images.ListImagesResponse
//	@Failure		400		{object}	map[string]interface{}	"Invalid request"
//	@Failure		422		{object}	ValidationErrorResponse		"Validation failed"
//	@Failure		500		{object}	map[string]interface{}	"Internal server error"
//	@Router			/v1/images/ [get]
func (h *ImageHandler) ListImages(c echo.Context) error {
	ctx := c.Request().Context()

	req := images.ListImagesRequest{}

	err := c.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	resp, err := h.imageUseCase.ListImages(ctx, &req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

// GetAllImagesRealtimeUpdates godoc
//
//	@Summary		Stream all images updates
//	@Description	Server-Sent Events stream for all images updates
//	@Tags			images
//	@Produce		text/event-stream
//	@Success		200	{object}	images.Image			"Stream of image updates"
//	@Failure		500	{object}	map[string]interface{}	"Streaming unsupported"
//	@Router			/v1/images/sse [get]
func (h *ImageHandler) GetAllImagesRealtimeUpdates(c echo.Context) error {
	ctx := context.Background()
	fluser, ok := c.Response().Writer.(http.Flusher)
	if !ok {
		slog.ErrorContext(ctx, "streaming unsupported by response writer")
		return echo.NewHTTPError(http.StatusInternalServerError, "Streaming unsupported")
	}

	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")

	imageUpdates, closeCallback, err := h.imageUseCase.GetAllImagesUpdates(ctx)
	if err != nil {
		return err
	}
	defer closeCallback()

	for {
		select {
		case <-c.Request().Context().Done():
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
			fluser.Flush()
		}
	}
}

// GetImageRealtimeUpdate godoc
//
//	@Summary		Stream single image updates
//	@Description	Server-Sent Events stream for a specific image updates
//	@Tags			images
//	@Produce		text/event-stream
//	@Param			id	path		string					true	"Image ID (UUID)"
//	@Success		200	{object}	images.Image			"Stream of image updates"
//	@Failure		404	{object}	map[string]interface{}	"Image not found"
//	@Failure		500	{object}	map[string]interface{}	"Streaming unsupported"
//	@Router			/v1/images/{id}/sse [get]
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
			return echo.NewHTTPError(http.StatusNotFound, "Image not found")
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

// GetImage godoc
//
//	@Summary		Get image by ID
//	@Description	Retrieve a single image by its ID
//	@Tags			images
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Image ID (UUID)"
//	@Success		200	{object}	images.Image
//	@Failure		404	{object}	map[string]interface{}	"Image not found"
//	@Failure		500	{object}	map[string]interface{}	"Internal server error"
//	@Router			/v1/images/{id} [get]
func (h *ImageHandler) GetImage(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")

	image, err := h.imageUseCase.GetImage(ctx, id)
	if err != nil {
		if errors.Is(err, ports.ErrImageNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Image not found")
		}
		return err
	}

	return c.JSON(http.StatusOK, image)
}

// CreateImage godoc
//
//	@Summary		Create a new image
//	@Description	Create a new image processing request
//	@Tags			images
//	@Accept			json
//	@Produce		json
//	@Param			request	body		images.CreateImageRequest	true	"Create image request"
//	@Success		201		{object}	images.CreateImageResponse
//	@Failure		400		{object}	map[string]interface{}		"Invalid request body"
//	@Failure		422		{object}	ValidationErrorResponse		"Validation failed"
//	@Failure		500		{object}	map[string]interface{}		"Internal server error"
//	@Router			/v1/images/ [post]
func (h *ImageHandler) CreateImage(c echo.Context) error {
	ctx := c.Request().Context()
	req := images.CreateImageRequest{}

	err := c.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	resp, err := h.imageUseCase.CreateImageRequest(ctx, &req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, resp)
}

// UpdateImage godoc
//
//	@Summary		Update an image
//	@Description	Update image transformation settings
//	@Tags			images
//	@Accept			json
//	@Produce		json
//	@Param			request	body		images.UpdateImageRequest	true	"Update image request"
//	@Success		200		{object}	map[string]string
//	@Failure		400		{object}	map[string]interface{}	"Invalid request body"
//	@Failure		404		{object}	map[string]interface{}	"Image not found"
//	@Failure		422		{object}	ValidationErrorResponse		"Validation failed"
//	@Failure		500		{object}	map[string]interface{}	"Internal server error"
//	@Router			/v1/images/ [put]
func (h *ImageHandler) UpdateImage(c echo.Context) error {
	ctx := c.Request().Context()
	req := images.UpdateImageRequest{}

	err := c.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	err = h.imageUseCase.UpdateImage(ctx, &req)
	if err != nil {
		if errors.Is(err, ports.ErrImageNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Image not found")
		}
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Image updated successfully"})
}

// DeleteImage godoc
//
//	@Summary		Delete an image
//	@Description	Delete an image by its ID
//	@Tags			images
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Image ID (UUID)"
//	@Success		200	{object}	map[string]string
//	@Failure		404	{object}	map[string]interface{}	"Image not found"
//	@Failure		500	{object}	map[string]interface{}	"Internal server error"
//	@Router			/v1/images/{id} [delete]
func (h *ImageHandler) DeleteImage(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")

	err := h.imageUseCase.DeleteImage(ctx, id)
	if err != nil {
		if errors.Is(err, ports.ErrImageNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Image not found")
		}
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Image deleted successfully"})
}
