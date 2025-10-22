package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	healthgo "github.com/hellofresh/health-go/v5"
	"github.com/taldoflemis/sora-henkan/internal/core/application"
	imageProcessor "github.com/taldoflemis/sora-henkan/internal/infra/adapter/image_processor"
	objectStorer "github.com/taldoflemis/sora-henkan/internal/infra/adapter/object_storer"
	"github.com/taldoflemis/sora-henkan/internal/infra/adapter/persistence/postgres"
	"github.com/taldoflemis/sora-henkan/internal/infra/handler/http"
	"github.com/taldoflemis/sora-henkan/internal/infra/telemetry"
	"github.com/taldoflemis/sora-henkan/settings"
)

type APISettings struct {
	App            settings.AppSettings            `mapstructure:"app" validate:"required"`
	Database       settings.DatabaseSettings       `mapstructure:"database" validate:"required"`
	OpenTelemetry  settings.OpenTelemetrySettings  `mapstructure:"opentelemetry" validate:"required"`
	HTTP           settings.HTTPSettings           `mapstructure:"http" validate:"required"`
	ImageProcessor settings.ImageProcessorSettings `mapstructure:"image-processor" validate:"required"`
	ObjectStorer   settings.ObjectStorerSettings   `mapstructure:"object-storer" validate:"required"`
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

	slog.InfoContext(ctx, "Launching app")

	slog.InfoContext(ctx, "Loading config")
	settings, err := settings.LoadConfig[APISettings]("API", settings.BaseSettings)
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

	router := http.NewRouter(&settings.HTTP, &settings.App)
	prefixedGroup := router.GetGroup()

	// Create adapters
	imageScaler := imageProcessor.NewVipsImageProcessor()
	objectStorerAdapter := objectStorer.NewMinioObjectStorer(minioClient)
	imageRepository := postgres.NewPostgresImageRepository(pgxpool)

	// Create usecases
	imageUseCase := application.NewImageUseCase(
		imageRepository,
		imageScaler,
		objectStorerAdapter,
		settings.ImageProcessor.BucketName,
	)

	// Register handlers
	healthHandler := http.NewHealthHandler(health)
	healthHandler.RegisterRoute(prefixedGroup)
	imageHandler := http.NewImageHandler(imageUseCase)
	imageHandler.RegisterRoute(prefixedGroup)

	errChan := make(chan error, 1)
	go func() {
		errChan <- router.Start()
	}()

	select {
	case err = <-errChan:
		slog.ErrorContext(ctx, "error when running server", slog.Any("err", err))
		return
	case <-ctx.Done():
		// Wait for first Signal arrives
	}

	err = router.Shutdown(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to shutdown gracefully the server", slog.Any("err", err))
		return
	}

	slog.InfoContext(ctx, "App stopped gracefully")
	retcode = 0
}
