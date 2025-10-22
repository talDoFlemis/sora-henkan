package application

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
	"github.com/taldoflemis/sora-henkan/internal/core/domain/images"
	"github.com/taldoflemis/sora-henkan/internal/core/ports"
	"github.com/taldoflemis/sora-henkan/internal/infra/telemetry"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("")

// Allowed MIME types for images
var allowedMIMETypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
}

var (
	rawImagePath         = "raw-images"
	transformedImagePath = "transformed-images"
)

type ImageUseCase struct {
	imageRepository ports.ImageRepository
	imageScaler     ports.ImageScaler
	objectStorer    ports.ObjectStorer
	imagesBucket    string
	publisher       message.Publisher
	subscriber      message.Subscriber
	imageTopic      string
}

func NewImageUseCase(
	publisher message.Publisher,
	subscriber message.Subscriber,
	imageRepository ports.ImageRepository,
	imageScaler ports.ImageScaler,
	objectStorer ports.ObjectStorer,
	imagesBucket string,
	imageTopic string,
) *ImageUseCase {
	return &ImageUseCase{
		publisher:       publisher,
		subscriber:      subscriber,
		imageRepository: imageRepository,
		imageScaler:     imageScaler,
		imagesBucket:    imagesBucket,
		objectStorer:    objectStorer,
		imageTopic:      imageTopic,
	}
}

func (u *ImageUseCase) CreateImageRequest(ctx context.Context, req *images.CreateImageRequest) (*images.CreateImageResponse, error) {
	ctx, span := tracer.Start(ctx, "ImageUseCase.CreateImage", trace.WithAttributes(
		attribute.String("image.original_url", req.ImageURL),
	))
	defer span.End()

	imageEntity := images.Image{
		ID:               uuid.New(),
		OriginalImageURL: req.ImageURL,
		CreatedAt:        time.Now(),
		Status:           "pending",
		UpdatedAt:        time.Now(),
	}

	err := u.imageRepository.CreateNewImage(ctx, &imageEntity)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create image", slog.Any("err", err))
		telemetry.RegisterSpanError(span, err)
		return nil, err
	}

	payload, err := json.Marshal(imageEntity)
	if err != nil {
		slog.ErrorContext(ctx, "failed to marshal image", slog.Any("err", err))
		telemetry.RegisterSpanError(span, err)
		return nil, err
	}

	message := message.NewMessageWithContext(ctx, watermill.NewUUID(), payload)

	err = u.publisher.Publish(u.imageTopic, message)
	if err != nil {
		slog.ErrorContext(ctx, "failed to publish image request", slog.Any("err", err))
		telemetry.RegisterSpanError(span, err)
		return nil, err
	}

	return &images.CreateImageResponse{
		ID: imageEntity.ID.String(),
	}, nil
}

func (u *ImageUseCase) GetImage(ctx context.Context, id string) (*images.Image, error) {
	ctx, span := tracer.Start(ctx, "ImageUseCase.GetImage", trace.WithAttributes(attribute.String("image.id", id)))
	defer span.End()

	storedImage, err := u.imageRepository.FindImageByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "failed to find image", slog.Any("err", err))
		telemetry.RegisterSpanError(span, err)
		return nil, err
	}

	return storedImage, nil
}

func (u *ImageUseCase) UpdateImage(ctx context.Context, req *images.UpdateImageRequest) error {
	ctx, span := tracer.Start(ctx, "ImageUseCase.UpdateImage", trace.WithAttributes(attribute.String("image.id", req.ID)))
	defer span.End()

	imageEntity, err := u.imageRepository.FindImageByID(ctx, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to find image", slog.Any("err", err))
		telemetry.RegisterSpanError(span, err)
		return err
	}

	imageEntity.ScaleTransformation = req.ScaleTransformation

	err = u.imageRepository.UpdateImage(ctx, imageEntity)
	if err != nil {
		slog.ErrorContext(ctx, "failed to update image", slog.Any("err", err))
		telemetry.RegisterSpanError(span, err)
		return err
	}

	payload, err := json.Marshal(imageEntity)
	if err != nil {
		slog.ErrorContext(ctx, "failed to marshal image", slog.Any("err", err))
		telemetry.RegisterSpanError(span, err)
		return err
	}

	message := message.NewMessageWithContext(ctx, watermill.NewUUID(), payload)

	err = u.publisher.Publish(u.imageTopic, message)
	if err != nil {
		slog.ErrorContext(ctx, "failed to publish image request", slog.Any("err", err))
		telemetry.RegisterSpanError(span, err)
		return err
	}

	return nil
}

func (u *ImageUseCase) DeleteImage(ctx context.Context, id string) error {
	ctx, span := tracer.Start(ctx, "ImageUseCase.DeleteImage", trace.WithAttributes(attribute.String("image.id", id)))
	defer span.End()

	err := u.imageRepository.DeleteImage(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "failed to delete image", slog.Any("err", err))
		telemetry.RegisterSpanError(span, err)
		return err
	}

	err = u.objectStorer.Delete(ctx, rawImagePath+"/"+id, u.imagesBucket)
	if err != nil {
		slog.ErrorContext(ctx, "failed to delete image", slog.Any("err", err))
		telemetry.RegisterSpanError(span, err)
		return err
	}

	err = u.objectStorer.Delete(ctx, transformedImagePath+"/"+id, u.imagesBucket)
	if err != nil {
		slog.ErrorContext(ctx, "failed to delete image", slog.Any("err", err))
		telemetry.RegisterSpanError(span, err)
		return err
	}

	slog.InfoContext(ctx, "image deleted successfully", slog.String("image_id", id))

	return nil
}

func (u *ImageUseCase) ProcessImage(ctx context.Context, req *images.ProcessImageRequest) error {
	ctx, span := tracer.Start(ctx, "ImageUseCase.ProcessImage", trace.WithAttributes(
		attribute.String("image.id", req.ID),
		attribute.String("image.original_url", req.OriginalImageURL),
		attribute.String("image.storage_key", req.StorageKey),
	))
	defer span.End()

	imageEntity, err := u.imageRepository.FindImageByID(ctx, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to find image", slog.Any("err", err))
		telemetry.RegisterSpanError(span, err)
		return err
	}

	if req.StorageKey == "" {
		slog.WarnContext(ctx, "image is not stored yet, fetching image")
		storageKey, err := u.fetchAndStoreImage(ctx, req)
		if err != nil {
			slog.ErrorContext(ctx, "failed to fetch and store image", slog.Any("err", err))
			telemetry.RegisterSpanError(span, err)
			return err
		}

		req.StorageKey = storageKey
	}

	imageData, err := u.objectStorer.Get(ctx, req.StorageKey, u.imagesBucket)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get image", slog.Any("err", err))
		telemetry.RegisterSpanError(span, err)
		return err
	}

	imageProcessed, err := u.imageScaler.Scale(ctx, imageData, 100, 100)
	if err != nil {
		slog.ErrorContext(ctx, "failed to scale image", slog.Any("err", err))
		telemetry.RegisterSpanError(span, err)
		return err
	}

	transformedImagePath := transformedImagePath + "/" + imageEntity.ID.String()

	err = u.objectStorer.Store(ctx, transformedImagePath, u.imagesBucket, imageEntity.MimeType, bytes.NewBuffer(imageProcessed))
	if err != nil {
		slog.ErrorContext(ctx, "failed to store transformed image", slog.Any("err", err))
		telemetry.RegisterSpanError(span, err)
		return err
	}

	imageEntity.Status = "processed"
	imageEntity.TransformedImageKey = transformedImagePath
	imageEntity.UpdatedAt = time.Now()

	err = u.imageRepository.UpdateImage(ctx, imageEntity)
	if err != nil {
		slog.ErrorContext(ctx, "failed to update image", slog.Any("err", err))
		telemetry.RegisterSpanError(span, err)
		return err
	}

	slog.InfoContext(ctx, "image processed successfully", slog.String("image_id", imageEntity.ID.String()))

	return nil
}

func (u *ImageUseCase) ListImages(ctx context.Context, req *images.ListImagesRequest) (*images.ListImagesResponse, error) {
	ctx, span := tracer.Start(ctx, "ImageUseCase.ListImages", trace.WithAttributes(
		attribute.Int64("page", int64(req.Page)),
		attribute.Int64("limit", int64(req.Limit)),
	))
	defer span.End()

	resp, err := u.imageRepository.FindAllImages(ctx, req)
	if err != nil {
		slog.ErrorContext(ctx, "failed to list images", slog.Any("err", err))
		telemetry.RegisterSpanError(span, err)
		return nil, err
	}

	return resp, nil
}

func (u *ImageUseCase) fetchAndStoreImage(ctx context.Context, req *images.ProcessImageRequest) (string, error) {
	ctx, span := tracer.Start(ctx, "ImageUseCase.fetchAndStoreImage", trace.WithAttributes(
		attribute.String("image.original_url", req.OriginalImageURL),
	))
	defer span.End()

	imageRequets, err := http.NewRequestWithContext(ctx, http.MethodGet, req.OriginalImageURL, nil)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create image request", slog.Any("err", err))
		telemetry.RegisterSpanError(span, err)
		return "", err
	}

	resp, err := http.DefaultClient.Do(imageRequets)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get image", slog.Any("err", err))
		telemetry.RegisterSpanError(span, err)
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("received non-OK HTTP status: %s", resp.Status)
		slog.ErrorContext(ctx, "received non-OK HTTP status", slog.Int("status_code", resp.StatusCode))
		telemetry.RegisterSpanError(span, err)
		return "", err
	}

	contentTypeHeader := resp.Header.Get("Content-Type")
	mimeType := strings.ToLower(strings.Split(contentTypeHeader, ";")[0])

	if !allowedMIMETypes[mimeType] {
		err = fmt.Errorf("disallowed MIME type: %s", mimeType)
		slog.ErrorContext(ctx, "disallowed MIME type", slog.String("mime_type", mimeType))
		telemetry.RegisterSpanError(span, err)
		return "", err
	}

	id := uuid.New()

	bodyData, err := io.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("failed to read response body: %w", err)
		slog.ErrorContext(ctx, "failed to read response body", slog.Any("err", err))
		telemetry.RegisterSpanError(span, err)
		return "", err
	}

	rawImageKey := rawImagePath + "/" + id.String()

	err = u.objectStorer.Store(ctx, rawImageKey, u.imagesBucket, mimeType, bytes.NewBuffer(bodyData))
	if err != nil {
		slog.ErrorContext(ctx, "failed to store image", slog.Any("err", err))
		telemetry.RegisterSpanError(span, err)
		return "", err
	}

	slog.InfoContext(ctx, "image stored successfully", slog.String("image_id", id.String()), slog.String("raw_image_key", rawImageKey))

	return rawImageKey, nil
}

func (u *ImageUseCase) GetImageRealtimeUpdate(ctx context.Context, id string) (chan *images.Image, func() error, error) {
	ctx, span := tracer.Start(ctx, "ImageUseCase.GetImageRealtimeUpdate", trace.WithAttributes(attribute.String("image.id", id)))
	defer span.End()

	_, err := u.imageRepository.FindImageByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "failed to find image", slog.Any("err", err))
		telemetry.RegisterSpanError(span, err)
		return nil, nil, err
	}

	ctx, cancel := context.WithCancel(ctx)

	ch, err := u.subscriber.Subscribe(ctx, u.imageTopic)
	if err != nil {
		slog.ErrorContext(ctx, "failed to subscribe to events", slog.Any("err", err))
		telemetry.RegisterSpanError(span, err)
		cancel()
		return nil, nil, err
	}

	imagesChan := make(chan *images.Image, len(ch))

	go func() {
		for msg := range ch {
			image := images.Image{}

			err := json.Unmarshal(msg.Payload, &image)
			if err != nil {
				slog.ErrorContext(ctx, "failed to unmarshal image", slog.Any("err", err))
				msg.Nack()
				continue
			}

			if image.ID.String() != id {
				msg.Nack()
				continue
			}

			imagesChan <- &image
		}
	}()

	return imagesChan, func() error {
		cancel()
		return nil
	}, nil
}

func (u *ImageUseCase) GetAllImagesUpdates(ctx context.Context) (chan *images.Image, func() error, error) {
	ctx, span := tracer.Start(ctx, "ImageUseCase.GetAllImagesUpdates")
	defer span.End()

	ctx, cancel := context.WithCancel(ctx)

	ch, err := u.subscriber.Subscribe(ctx, u.imageTopic)
	if err != nil {
		slog.ErrorContext(ctx, "failed to subscribe to events", slog.Any("err", err))
		telemetry.RegisterSpanError(span, err)
		cancel()
		return nil, nil, err
	}

	imagesChan := make(chan *images.Image, len(ch))

	go func() {
		for msg := range ch {
			image := images.Image{}

			err := json.Unmarshal(msg.Payload, &image)
			if err != nil {
				slog.ErrorContext(ctx, "failed to unmarshal image", slog.Any("err", err))
				msg.Nack()
				continue
			}

			imagesChan <- &image
		}
	}()

	return imagesChan, func() error {
		cancel()
		return nil
	}, nil
}
