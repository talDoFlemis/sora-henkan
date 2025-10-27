package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	wotelfloss "github.com/dentech-floss/watermill-opentelemetry-go-extra/pkg/opentelemetry"
	healthgo "github.com/hellofresh/health-go/v5"
	"github.com/labstack/echo/v4"
	"github.com/taldoflemis/sora-henkan/internal/core/application"
	"github.com/taldoflemis/sora-henkan/internal/core/domain/images"
	imageProcessor "github.com/taldoflemis/sora-henkan/internal/infra/adapter/image_processor"
	objectStorer "github.com/taldoflemis/sora-henkan/internal/infra/adapter/object_storer"
	"github.com/taldoflemis/sora-henkan/internal/infra/adapter/persistence/postgres"
	"github.com/taldoflemis/sora-henkan/internal/infra/telemetry"
	"github.com/taldoflemis/sora-henkan/settings"
	wotel "github.com/voi-oss/watermill-opentelemetry/pkg/opentelemetry"
)

type WorkerSettings struct {
	App            settings.AppSettings            `mapstructure:"app" validate:"required"`
	Database       settings.DatabaseSettings       `mapstructure:"database" validate:"required"`
	OpenTelemetry  settings.OpenTelemetrySettings  `mapstructure:"opentelemetry" validate:"required"`
	HTTP           settings.HTTPSettings           `mapstructure:"http" validate:"required"`
	ImageProcessor settings.ImageProcessorSettings `mapstructure:"image-processor" validate:"required"`
	ObjectStorer   settings.ObjectStorerSettings   `mapstructure:"object-storer" validate:"required"`
	Watermill      settings.WatermillSettings      `mapstructure:"watermill" validate:"required"`
}

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()
	retcode := 1

	defer func() {
		os.Exit(retcode)
	}()

	slog.InfoContext(ctx, "Launching worker")

	slog.InfoContext(ctx, "Loading config")
	settings, err := settings.LoadConfig[WorkerSettings]("WORKER", settings.BaseSettings)
	if err != nil {
		slog.ErrorContext(ctx, "failed to load config", slog.Any("err", err))
		return
	}

	slog.InfoContext(ctx, "Setting up opentelemetry")
	otelShutdown, err := telemetry.SetupOTelSDK(ctx, settings.App, settings.OpenTelemetry)
	if err != nil {
		slog.Error("failed to setup telemetry", slog.Any("err", err))
		return
	}

	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
		if err != nil {
			slog.ErrorContext(
				ctx,
				"failed to shutdown opentelemetry providers",
				slog.Any("err", err),
			)
		}
	}()

	slog.InfoContext(ctx, "Initializing PostgreSQL client")
	pgxpool, err := postgres.NewPool(ctx, settings.Database)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create postgres pool", slog.Any("err", err))
		return
	}

	slog.InfoContext(ctx, "Initializing MinIO client")
	minioClient, err := settings.ObjectStorer.NewMinioClient(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to initialize MinIO client", slog.Any("err", err))
		return
	}

	health, err := healthgo.New(
		healthgo.WithComponent(healthgo.Component{
			Name:    settings.App.Name,
			Version: settings.App.Version,
		}),
		healthgo.WithChecks(
			healthgo.Config{
				Name:  "postgres",
				Check: pgxpool.Ping,
			},
			healthgo.Config{
				Name: "s3client",
				Check: func(ctx context.Context) error {
					_, err := minioClient.ListBuckets(ctx)
					return err
				},
			},
		),
	)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create health checker", slog.Any("err", err))
		return
	}

	// Start health check HTTP server
	e := echo.New()
	e.HidePort = true
	e.HideBanner = true

	e.GET("/healthz", func(c echo.Context) error {
		check := health.Measure(c.Request().Context())

		statusCode := http.StatusOK
		if check.Status != healthgo.StatusOK {
			statusCode = http.StatusServiceUnavailable
		}

		return c.JSON(statusCode, check)
	})

	healthServer := &http.Server{
		Addr:    settings.HTTP.IP + ":" + settings.HTTP.Port,
		Handler: e,
	}

	go func() {
		slog.InfoContext(ctx, "Starting health check server", slog.String("addr", healthServer.Addr))
		if err := healthServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.ErrorContext(ctx, "health check server failed", slog.Any("err", err))
		}
	}()

	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := healthServer.Shutdown(shutdownCtx); err != nil {
			slog.ErrorContext(ctx, "failed to shutdown health check server", slog.Any("err", err))
		}
	}()

	slog.InfoContext(ctx, "Setting up Watermill")

	publisher, err := settings.Watermill.Broker.NewPublisher()
	if err != nil {
		slog.ErrorContext(ctx, "failed to create publisher", slog.Any("err", err))
		return
	}

	subscriber, err := settings.Watermill.Broker.NewSubscriber()
	if err != nil {
		slog.ErrorContext(ctx, "failed to create subscriber", slog.Any("err", err))
		return
	}

	// Create adapters
	transformerFactory := imageProcessor.NewTransformerFactory(
		imageProcessor.NewVipsResizeTransformer(),
		imageProcessor.NewVipsGrayscaleTransformer(),
		imageProcessor.NewVipsTrimTransformer(),
		imageProcessor.NewVipsBlurTransformer(),
		imageProcessor.NewVipsRotateTransformer(),
	)
	pipelineProcessor := imageProcessor.NewPipeline(transformerFactory)
	objectStorerAdapter := objectStorer.NewMinioObjectStorer(minioClient)
	imageRepository := postgres.NewPostgresImageRepository(pgxpool)

	// Create usecases
	imageUseCase := application.NewImageUseCase(
		publisher,
		subscriber,
		imageRepository,
		pipelineProcessor,
		objectStorerAdapter,
		settings.ImageProcessor.BucketName,
		settings.Watermill.ImageTopic,
	)

	slog.InfoContext(ctx, "Setting up Watermill router")

	// Create Watermill router
	router, err := message.NewRouter(message.RouterConfig{}, nil)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create watermill router", slog.Any("err", err))
		return
	}

	// Add handler for pending images
	router.AddConsumerHandler(
		"process_pending_images",
		settings.Watermill.ImageTopic,
		subscriber,
		func(msg *message.Message) error {
			return handlePendingImage(msg, imageUseCase)
		},
	)

	router.AddMiddleware(wotelfloss.ExtractRemoteParentSpanContext())
	router.AddMiddleware(wotel.Trace())
	router.AddMiddleware(middleware.Recoverer)
	router.AddMiddleware(middleware.Retry{
		MaxRetries: 5,
	}.Middleware)
	router.AddPlugin(plugin.SignalsHandler)

	slog.InfoContext(ctx, "Starting Watermill router")

	errChan := make(chan error, 1)
	go func() {
		errChan <- router.Run(ctx)
	}()

	select {
	case err = <-errChan:
		if err != nil {
			slog.ErrorContext(ctx, "error when running watermill router", slog.Any("err", err))
			return
		}
	case <-ctx.Done():
		// Wait for first Signal arrives
		slog.InfoContext(ctx, "Shutdown signal received")
	}

	slog.InfoContext(ctx, "Closing Watermill router")
	if err := router.Close(); err != nil {
		slog.ErrorContext(ctx, "failed to close watermill router gracefully", slog.Any("err", err))
		return
	}

	slog.InfoContext(ctx, "Worker stopped gracefully")
	retcode = 0
}

// handlePendingImage processes pending image messages
func handlePendingImage(msg *message.Message, imageUseCase *application.ImageUseCase) error {
	ctx := msg.Context()

	slog.InfoContext(ctx, "Processing pending image message", slog.String("message_id", msg.UUID))

	var processReq images.ProcessImageRequest
	if err := json.Unmarshal(msg.Payload, &processReq); err != nil {
		slog.ErrorContext(ctx, "failed to unmarshal process image request from message", slog.Any("err", err))
		// Nack the message - it will be redelivered or sent to DLQ depending on configuration
		return err
	}

	slog.InfoContext(ctx, "Process image request received",
		slog.String("image_id", processReq.ID),
		slog.String("original_url", processReq.OriginalImageURL),
		slog.String("storage_key", processReq.StorageKey),
	)

	// Process the image
	if err := imageUseCase.ProcessImage(ctx, &processReq); err != nil {
		slog.ErrorContext(ctx, "failed to process image",
			slog.String("image_id", processReq.ID),
			slog.Any("err", err),
		)
		// Return error to nack the message
		return err
	}

	slog.InfoContext(ctx, "Image processed successfully",
		slog.String("image_id", processReq.ID),
		slog.String("message_id", msg.UUID),
	)

	msg.Ack()

	return nil
}
