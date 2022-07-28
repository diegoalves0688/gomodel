//go:build unit
// +build unit

package message

import (
	"net/http"
	"testing"

	"github.com/diegoalves0688/gomodel/internal/domain"
	"github.com/diegoalves0688/gomodel/mocks"
	"github.com/gavv/httpexpect/v2"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
)

func TestMessageHTTPHandler(t *testing.T) {
	t.Run("Returns a list of messages with success", func(t *testing.T) {
		// Arrange
		e := echo.New()
		usecaseMock := &mocks.MessageUseCase{}
		usecaseMock.On("List", mock.AnythingOfType("*context.emptyCtx")).Return(
			[]domain.Message{
				{
					ID:       uuid.MustParse("08dcd11a-2231-4a53-8db2-bde5fb2dc7e9"),
					Receiver: "Paulo",
					Sender:   "Maria",
					Content:  "any",
				},
			},
			nil,
		)
		h := NewMessageHandler(usecaseMock)
		e.GET("/messages", h.List)

		// Act
		ex := HTTPExpect(t, e).
			GET("/messages")

		// Assert
		ex.Expect().Status(http.StatusOK).JSON()
	})
}

func HTTPExpect(t *testing.T, h http.Handler) *httpexpect.Expect {
	return httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(h),
			Jar:       httpexpect.NewJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})
}
