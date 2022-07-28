package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func DefaultHTTPErrorHandler(logger *zap.Logger) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
		}
		logger.Error("Error", zap.Error(err))
		if err := c.NoContent(code); err != nil {
			logger.Error("Error to return", zap.Error(err))
		}
	}
}
