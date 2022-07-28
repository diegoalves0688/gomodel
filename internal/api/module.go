package api

import (
	"github.com/diegoalves0688/gomodel/internal/api/message"
	"github.com/diegoalves0688/gomodel/internal/infra/repository"
	"github.com/diegoalves0688/gomodel/internal/usecase"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

func routes(e *echo.Echo, h message.MessageHandler) {
	e.GET("/messages", h.List)
	e.POST("/messages", h.Create)
	e.GET("/messages/:id", h.FindByID)
}

var Module = fx.Options(
	fx.Provide(
		message.NewMessageHandler,
		usecase.NewMessageUseCaseImpl,
		repository.NewMessageRepositoryImpl,
	),
	fx.Invoke(routes),
)
