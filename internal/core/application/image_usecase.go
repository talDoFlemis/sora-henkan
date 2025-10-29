package application

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime"
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
	imageRepository   ports.ImageRepository
	pipelineProcessor ports.ImagePipelineProcessor
	objectStorer      ports.ObjectStorer
	imagesBucket      string
	publisher         message.Publisher
	subscriber        message.Subscriber
	imageTopic        string
}

func NewImageUseCase(
	publisher message.Publisher,
	subscriber message.Subscriber,
	imageRepository ports.ImageRepository,
	pipelineProcessor ports.ImagePipelineProcessor,
	objectStorer ports.ObjectStorer,
	imagesBucket string,
	imageTopic string,
) *ImageUseCase {
	return &ImageUseCase{
		publisher:         publisher,
		subscriber:        subscriber,
		imageRepository:   imageRepository,
		pipelineProcessor: pipelineProcessor,
		imagesBucket:      imagesBucket,
		objectStorer:      objectStorer,
		imageTopic:        imageTopic,
	}
}

func (u *ImageUseCase) CreateImageRequest(ctx context.Context, req *images.CreateImageRequest) (*images.CreateImageResponse, error) {
	ctx, span := tracer.Start(ctx, "ImageUseCase.CreateImage", trace.WithAttributes(
		attribute.String("image.original_url", req.ImageURL),
	))
	defer span.End()

	// Validate request
	if err := ValidateStruct(req); err != nil {
		slog.ErrorContext(ctx, "validation failed", slog.Any("err", err))
		telemetry.RegisterSpanError(span, err)
		return nil, err
	}

	// Validate transformations
	if err := u.pipelineProcessor.ValidateTransformations(ctx, req.Transformations); err != nil {
		slog.ErrorContext(ctx, "transformation validation failed", slog.Any("err", err))
		telemetry.RegisterSpanError(span, err)
		return nil, err
	}

	imageEntity := images.Image{
		ID:               uuid.New(),
		OriginalImageURL: req.ImageURL,
		Transformations:  images.TransformationList(req.Transformations),
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

	processReq := &images.ProcessImageRequest{
		ID:               imageEntity.ID.String(),
		OriginalImageURL: imageEntity.OriginalImageURL,
		StorageKey:       imageEntity.ObjectStorageImageKey,
		Transformations:  req.Transformations,
	}

	payload, err := json.Marshal(processReq)
	if err != nil {
		slog.ErrorContext(ctx, "failed to marshal process image request", slog.Any("err", err))
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

func (u *ImageUseCase) ValidateTransformations(ctx context.Context, transformations []images.TransformationRequest) error {
	ctx, span := tracer.Start(ctx, "ImageUseCase.ValidateTransformations")
	defer span.End()

	err := u.pipelineProcessor.ValidateTransformations(ctx, transformations)
	if err != nil {
		slog.ErrorContext(ctx, "transformation validation failed", slog.Any("err", err))
		telemetry.RegisterSpanError(span, err)
		return err
	}

	return nil
}

func (u *ImageUseCase) UpdateImage(ctx context.Context, req *images.UpdateImageRequest) error {
	ctx, span := tracer.Start(ctx, "ImageUseCase.UpdateImage", trace.WithAttributes(attribute.String("image.id", req.ID)))
	defer span.End()

	// Validate request
	if err := ValidateStruct(req); err != nil {
		slog.ErrorContext(ctx, "validation failed", slog.Any("err", err))
		telemetry.RegisterSpanError(span, err)
		return err
	}

	// Validate transformations
	if err := u.pipelineProcessor.ValidateTransformations(ctx, req.Transformations); err != nil {
		slog.ErrorContext(ctx, "transformation validation failed", slog.Any("err", err))
		telemetry.RegisterSpanError(span, err)
		return err
	}

	imageEntity, err := u.imageRepository.FindImageByID(ctx, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to find image", slog.Any("err", err))
		telemetry.RegisterSpanError(span, err)
		return err
	}

	imageEntity.Transformations = images.TransformationList(req.Transformations)
	imageEntity.UpdatedAt = time.Now()

	err = u.imageRepository.UpdateImage(ctx, imageEntity)
	if err != nil {
		slog.ErrorContext(ctx, "failed to update image", slog.Any("err", err))
		telemetry.RegisterSpanError(span, err)
		return err
	}

	processReq := &images.ProcessImageRequest{
		ID:               imageEntity.ID.String(),
		OriginalImageURL: imageEntity.OriginalImageURL,
		StorageKey:       imageEntity.ObjectStorageImageKey,
		Transformations:  req.Transformations,
	}

	payload, err := json.Marshal(processReq)
	if err != nil {
		slog.ErrorContext(ctx, "failed to marshal process image request", slog.Any("err", err))
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

	// Validate request
	if err := ValidateStruct(req); err != nil {
		slog.ErrorContext(ctx, "validation failed", slog.Any("err", err))
		telemetry.RegisterSpanError(span, err)
		return err
	}

	// Validate transformations
	if err := u.pipelineProcessor.ValidateTransformations(ctx, req.Transformations); err != nil {
		slog.ErrorContext(ctx, "transformation validation failed", slog.Any("err", err))
		telemetry.RegisterSpanError(span, err)
		return err
	}

	imageEntity, err := u.imageRepository.FindImageByID(ctx, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to find image", slog.Any("err", err))
		telemetry.RegisterSpanError(span, err)
		return err
	}

	if req.StorageKey == "" {
		slog.WarnContext(ctx, "image is not stored yet, fetching image")
		imageData, err := u.fetchAndStoreImage(ctx, req)
		if err != nil {
			slog.ErrorContext(ctx, "failed to fetch and store image", slog.Any("err", err))
			telemetry.RegisterSpanError(span, err)

			// Check if this is a non-retryable error (4xx status code)
			var nonRetryableErr *images.NonRetryableError
			if errors.As(err, &nonRetryableErr) {
				// Update image status to failed
				imageEntity.Status = "failed"
				imageEntity.ErrorMessage = nonRetryableErr.Error()
				imageEntity.UpdatedAt = time.Now()

				if updateErr := u.imageRepository.UpdateImage(ctx, imageEntity); updateErr != nil {
					slog.ErrorContext(ctx, "failed to update image status to failed", slog.Any("err", updateErr))
					// Return the original non-retryable error
				}

				// Return the non-retryable error so the worker can ACK the message
				return err
			}

			return err
		}

		req.StorageKey = imageData.ObjectStorageImageKey
		imageEntity.MimeType = imageData.MimeType
		imageEntity.ObjectStorageImageKey = imageData.ObjectStorageImageKey

		// Update image entity with storage information
		err = u.imageRepository.UpdateImage(ctx, imageEntity)
		if err != nil {
			slog.ErrorContext(ctx, "failed to update image with storage info", slog.Any("err", err))
			telemetry.RegisterSpanError(span, err)
			return err
		}
	}

	imageData, err := u.objectStorer.Get(ctx, req.StorageKey, u.imagesBucket)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get image", slog.Any("err", err), slog.String("bucket-name", u.imagesBucket))
		telemetry.RegisterSpanError(span, err)
		return err
	}

	var imageProcessed []byte

	processedBytes, err := u.pipelineProcessor.ProcessPipeline(ctx, imageData, req.Transformations)
	if err != nil {
		slog.ErrorContext(ctx, "failed to process image transformations", slog.Any("err", err))
		telemetry.RegisterSpanError(span, err)
		return err
	}

	slog.InfoContext(ctx, "image processed successfully", slog.String("image_id", req.ID), slog.Int("transformations", len(req.Transformations)))
	imageProcessed = processedBytes

	extensions, err := mime.ExtensionsByType(imageEntity.MimeType)
	if err != nil || len(extensions) == 0 {
		err = fmt.Errorf("failed to get file extension for MIME type: %s", imageEntity.MimeType)
		slog.ErrorContext(ctx, "failed to get file extension", slog.String("mime_type", imageEntity.MimeType))
		telemetry.RegisterSpanError(span, err)
		return err
	}

	transformedImagePath := transformedImagePath + "/" + imageEntity.ID.String() + extensions[0]

	err = u.objectStorer.Store(ctx, transformedImagePath, u.imagesBucket, imageEntity.MimeType, bytes.NewReader(imageProcessed))
	if err != nil {
		slog.ErrorContext(ctx, "failed to store transformed image", slog.Any("err", err), slog.String("bucket-name", u.imagesBucket))
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

	// Validate request
	if err := ValidateStruct(req); err != nil {
		slog.ErrorContext(ctx, "validation failed", slog.Any("err", err))
		telemetry.RegisterSpanError(span, err)
		return nil, err
	}

	resp, err := u.imageRepository.FindAllImages(ctx, req)
	if err != nil {
		slog.ErrorContext(ctx, "failed to list images", slog.Any("err", err))
		telemetry.RegisterSpanError(span, err)
		return nil, err
	}

	return resp, nil
}

func (u *ImageUseCase) fetchAndStoreImage(ctx context.Context, req *images.ProcessImageRequest) (*images.Image, error) {
	ctx, span := tracer.Start(ctx, "ImageUseCase.fetchAndStoreImage", trace.WithAttributes(
		attribute.String("image.original_url", req.OriginalImageURL),
	))
	defer span.End()

	imageRequets, err := http.NewRequestWithContext(ctx, http.MethodGet, req.OriginalImageURL, nil)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create image request", slog.Any("err", err))
		telemetry.RegisterSpanError(span, err)
		return nil, err
	}

	resp, err := http.DefaultClient.Do(imageRequets)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get image", slog.Any("err", err))
		telemetry.RegisterSpanError(span, err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// If status code is between 400 and 500 (client errors), mark as failed and don't retry
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			err = fmt.Errorf("client error HTTP status: %s", resp.Status)
			slog.ErrorContext(ctx, "received client error HTTP status", slog.Int("status_code", resp.StatusCode))
			telemetry.RegisterSpanError(span, err)

			// Return a special error type that indicates this should not be retried
			return nil, images.NewNonRetryableError(err)
		}

		// For other errors (5xx), return a regular error that will be retried
		err = fmt.Errorf("received non-OK HTTP status: %s", resp.Status)
		slog.ErrorContext(ctx, "received non-OK HTTP status", slog.Int("status_code", resp.StatusCode))
		telemetry.RegisterSpanError(span, err)
		return nil, err
	}

	contentTypeHeader := resp.Header.Get("Content-Type")
	mimeType := strings.ToLower(strings.Split(contentTypeHeader, ";")[0])

	if !allowedMIMETypes[mimeType] {
		err = fmt.Errorf("disallowed MIME type: %s", mimeType)
		slog.ErrorContext(ctx, "disallowed MIME type", slog.String("mime_type", mimeType))
		telemetry.RegisterSpanError(span, err)
		return nil, err
	}

	bodyData, err := io.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("failed to read response body: %w", err)
		slog.ErrorContext(ctx, "failed to read response body", slog.Any("err", err))
		telemetry.RegisterSpanError(span, err)
		return nil, err
	}
	extensions, err := mime.ExtensionsByType(mimeType)
	if err != nil || len(extensions) == 0 {
		err = fmt.Errorf("failed to get file extension for MIME type: %s", mimeType)
		slog.ErrorContext(ctx, "failed to get file extension", slog.String("mime_type", mimeType))
		telemetry.RegisterSpanError(span, err)
		return nil, err
	}

	rawImageKey := rawImagePath + "/" + req.ID + extensions[0]

	err = u.objectStorer.Store(ctx, rawImageKey, u.imagesBucket, mimeType, bytes.NewBuffer(bodyData))
	if err != nil {
		slog.ErrorContext(ctx, "failed to store image", slog.Any("err", err), slog.String("bucket-name", u.imagesBucket))
		telemetry.RegisterSpanError(span, err)
		return nil, err
	}

	slog.InfoContext(ctx, "image stored successfully", slog.String("image_id", req.ID), slog.String("raw_image_key", rawImageKey))

	imageData := &images.Image{
		ObjectStorageImageKey: rawImageKey,
		MimeType:              mimeType,
	}

	return imageData, nil
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
