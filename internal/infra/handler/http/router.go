package http

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"
	"github.com/taldoflemis/sora-henkan/settings"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var tracer = otel.Tracer("")

type Router struct {
	e            *echo.Echo
	httpSettings *settings.HTTPSettings
	appSettings  *settings.AppSettings
}

func NewRouter(httpSettings *settings.HTTPSettings, appSettings *settings.AppSettings, dynamoClient *dynamodb.Client, dynamoTableName string) *Router {
	logger := slog.Default()
	e := echo.New()

	e.HideBanner = true

	// Set custom error handler
	e.HTTPErrorHandler = GlobalErrorHandler
	e.IPExtractor = CloudFlareExtractClientIPfunc
	e.Use(
		middleware.TimeoutWithConfig(middleware.TimeoutConfig{
			Timeout: time.Duration(httpSettings.Timeout) * time.Second, Skipper: SSETimeoutSkipper,
		}),
	)
	e.Use(middleware.Recover())

	e.Use(otelecho.Middleware(appSettings.Name,
		otelecho.WithMetricAttributeFn(func(r *http.Request) []attribute.KeyValue {
			return []attribute.KeyValue{
				attribute.String("client.ip", r.RemoteAddr),
				attribute.String("user.agent", r.UserAgent()),
			}
		}),
		otelecho.WithEchoMetricAttributeFn(func(c echo.Context) []attribute.KeyValue {
			return []attribute.KeyValue{
				attribute.String("handler.path", c.Path()),
				attribute.String("handler.method", c.Request().Method),
			}
		}),
	))
	e.Use(slogecho.New(logger))
	e.Use(DynamoDBAuditLogger(dynamoClient, dynamoTableName))
	e.Use(ValidationErrorMiddleware())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: httpSettings.CORS.Origins,
		AllowMethods: httpSettings.CORS.Methods,
		AllowHeaders: httpSettings.CORS.Headers,
	}))

	return &Router{
		e:            e,
		appSettings:  appSettings,
		httpSettings: httpSettings,
	}
}

func (r *Router) GetGroup() *echo.Group {
	return r.e.Group(r.httpSettings.Prefix)
}

func (r *Router) RegisterSwagger() {
	if r.httpSettings.SwaggerUIEnabled {
		slog.Info("Swagger UI enabled, registering swagger handler")
		swaggerHandler := NewSwaggerHandler()
		swaggerHandler.RegisterRoute(r.e)
	}
}

func (r *Router) Start() error {
	slog.Info("listening for requests", slog.String("ip", r.httpSettings.IP), slog.String("port", r.httpSettings.Port))
	address := r.httpSettings.IP + ":" + r.httpSettings.Port
	return r.e.Start(address)
}

func (r *Router) Shutdown(ctx context.Context) error {
	slog.Info("shutting down http server")
	return r.e.Shutdown(ctx)
}
