package main

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/diegoalves0688/gomodel/internal/api"
	"github.com/diegoalves0688/gomodel/pkg/config"
	"github.com/diegoalves0688/gomodel/pkg/db"
	"github.com/diegoalves0688/gomodel/pkg/logger"
	"github.com/diegoalves0688/gomodel/pkg/migrations"
	tracer "github.com/diegoalves0688/gomodel/pkg/tracer/datadog"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func newEcho(lc fx.Lifecycle, logger *zap.Logger) *echo.Echo {
	e := echo.New()

	// enable HTTP Requests tracing
	e.Use(tracer.NewEchoMiddleware())

	// enable recover from panic
	rcvCfg := middleware.RecoverConfig{
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			logger.Error("HTTP Server recovered from panic", zap.Error(err))
			return err
		},
	}
	e.Use(middleware.RecoverWithConfig(rcvCfg))

	e.HTTPErrorHandler = api.DefaultHTTPErrorHandler(logger)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	return e
	ç
}

func runTracer(lc fx.Lifecycle, cfg config.Config) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			tracer.Start(cfg)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			tracer.Stop()
			return nil
		},
	})
}

func runServer(lc fx.Lifecycle, e *echo.Echo, logger *zap.Logger) *echo.Echo {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			logger.Info("Starting HTTP server")
			go func() {
				if err := e.Start(":1323"); !errors.Is(err, http.ErrServerClosed) {
					logger.Fatal("Error running server", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping HTTP server")
			cancelCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()
			if err := e.Shutdown(cancelCtx); err != nil {
				logger.Error("Error shutting down server", zap.Error(err))
			} else {
				logger.Info("Server shutdown gracefully")
			}
			return logger.Sync()
		},
	})
	return e
}

// @title        Swagger API
// @version      1.0
// @description  The swagger doc for the Go Sample API.
func main() {
	config.InitProfileConfig("./", "yaml")

	fx.New(
		fx.Provide(
			config.LoadConfig,
			logger.NewZapLogger,
			db.ConnectPostgres,
			newEcho,
		),
		api.Module,
		fx.Invoke(migrations.RunPostgresMigrations("file://migrations")),
		fx.Invoke(runTracer),
		fx.Invoke(runServer),
	).Run()
}
