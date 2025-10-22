package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	healthgo "github.com/hellofresh/health-go/v5"
	"github.com/taldoflemis/sora-henkan/internal/infra/handler/http"
	"github.com/taldoflemis/sora-henkan/internal/infra/telemetry"
	"github.com/taldoflemis/sora-henkan/settings"
)

type APISettings struct {
	App           settings.AppSettings           `mapstructure:"app" validate:"required"`
	OpenTelemetry settings.OpenTelemetrySettings `mapstructure:"opentelemetry" validate:"required"`
	HTTP          settings.HTTPSettings          `mapstructure:"http" validate:"required"`
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

	health, err := healthgo.New(
		healthgo.WithComponent(healthgo.Component{
			Name:    settings.App.Name,
			Version: settings.App.Version,
		}),
	)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create health checker", slog.Any("err", err))
		return
	}

	router := http.NewRouter(&settings.HTTP, &settings.App)
	prefixedGroup := router.GetGroup()

	// Register handlers
	healthHandler := http.NewHealthHandler(health)
	healthHandler.RegisterRoute(prefixedGroup)

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
