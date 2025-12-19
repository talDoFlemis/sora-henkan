package http

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/taldoflemis/sora-henkan/internal/core/application"
	"github.com/taldoflemis/sora-henkan/internal/core/domain/images"
	"github.com/taldoflemis/sora-henkan/internal/core/ports"
	"github.com/taldoflemis/sora-henkan/pkg/http/api"
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
	imageHandlerGroup.GET("/:id/metadata", h.GetImageMetadata)
	imageHandlerGroup.POST("/", h.CreateImage)
	imageHandlerGroup.PUT("/", h.UpdateImage)
	imageHandlerGroup.DELETE("/:id", h.DeleteImage)
}

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

	// Convert domain images to API images
	apiImages := make([]api.Image, 0, len(resp.Data))
	for _, domainImage := range resp.Data {
		apiImage, err := api.ConvertDomainImageToAPI(&domainImage)
		if err != nil {
			slog.ErrorContext(ctx, "failed to convert image", slog.String("error", err.Error()))
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to convert images")
		}
		apiImages = append(apiImages, *apiImage)
	}

	apiResp := api.ListImagesResponse{
		Page:  resp.Page,
		Limit: resp.Limit,
		Count: resp.Count,
		Data:  apiImages,
	}

	return c.JSON(http.StatusOK, apiResp)
}

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
			// Convert domain image to API image
			apiImage, err := api.ConvertDomainImageToAPI(update)
			if err != nil {
				slog.ErrorContext(ctx, "failed to convert image for SSE", slog.String("error", err.Error()))
				continue
			}

			data, err := json.Marshal(apiImage)
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
			// Convert domain image to API image
			apiImage, err := api.ConvertDomainImageToAPI(update)
			if err != nil {
				slog.ErrorContext(ctx, "failed to convert image for SSE", slog.String("error", err.Error()))
				continue
			}

			data, err := json.Marshal(apiImage)
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
			return echo.NewHTTPError(http.StatusNotFound, "Image not found")
		}
		return err
	}

	// Convert domain image to API image
	apiImage, err := api.ConvertDomainImageToAPI(image)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to convert image")
	}

	return c.JSON(http.StatusOK, apiImage)
}

func (h *ImageHandler) GetImageMetadata(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	metadata, err := h.imageUseCase.GetImageMetadata(ctx, id)
	if err != nil {
		if errors.Is(err, ports.ErrMetadataNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Metadata not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get metadata")
	}

	return c.JSON(http.StatusOK, metadata)
}

func (h *ImageHandler) CreateImage(c echo.Context) error {
	ctx := c.Request().Context()
	req := api.CreateImageRequest{}

	err := c.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Convert API request to domain request
	domainReq, err := api.ConvertAPICreateImageRequestToDomain(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	resp, err := h.imageUseCase.CreateImageRequest(ctx, domainReq)
	if err != nil {
		return err
	}

	// Convert response - parse string UUID to uuid.UUID
	id, err := uuid.Parse(resp.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to parse image ID")
	}

	apiResp := api.CreateImageResponse{
		Id: id,
	}

	return c.JSON(http.StatusCreated, apiResp)
}

func (h *ImageHandler) UpdateImage(c echo.Context) error {
	ctx := c.Request().Context()
	req := api.UpdateImageRequest{}

	err := c.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Convert API request to domain request
	domainReq, err := api.ConvertAPIUpdateImageRequestToDomain(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = h.imageUseCase.UpdateImage(ctx, domainReq)
	if err != nil {
		if errors.Is(err, ports.ErrImageNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Image not found")
		}
		return err
	}

	message := "Image updated successfully"
	return c.JSON(http.StatusOK, api.MessageResponse{Message: &message})
}

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

	message := "Image deleted successfully"
	return c.JSON(http.StatusOK, api.MessageResponse{Message: &message})
}
