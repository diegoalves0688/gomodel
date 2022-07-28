package message

import (
	"net/http"

	"github.com/diegoalves0688/gomodel/internal/usecase"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type MessageHandler struct {
	Usecase usecase.MessageUseCase
}

func NewMessageHandler(u usecase.MessageUseCase) MessageHandler {
	return MessageHandler{u}
}

// List godoc
// @Description  List messages
// @Tags         messages
// @Produce      json
// @Success      200  {object}  []MessageDTO
// @Router       /messages [get].
func (h *MessageHandler) List(c echo.Context) error {
	rows, err := h.Usecase.List(c.Request().Context())
	if err != nil {
		return err
	}

	var out []MessageDTO
	for i := range rows {
		if m, err := ToMessageDTO(&rows[i]); err == nil {
			out = append(out, m)
		} else {
			return err
		}
	}

	return c.JSON(http.StatusOK, out)
}

// Create godoc
// @Description  Create messages
// @Tags         messages
// @Accept       json
// @Produce      json
// @Param        message  body  CreateMessageDTO  true  "Message input"
// @Success      200
// @Router       /messages [post].
func (h *MessageHandler) Create(c echo.Context) error {
	in := new(CreateMessageDTO)

	if err := c.Bind(in); err != nil {
		return err
	}

	entity, err := in.ToEntity()
	if err != nil {
		return err
	}

	err = h.Usecase.Create(c.Request().Context(), &entity)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

// FindByID godoc
// @Description  FindByID message by id
// @Tags         messages
// @Produce      json
// @Success      200  {object}  MessageDTO
// @Param        id   path      string  true  "Message ID"
// @Router       /messages/{id} [get].
func (h *MessageHandler) FindByID(c echo.Context) error {
	var in string
	var id uuid.UUID

	if err := echo.PathParamsBinder(c).String("id", &in).BindError(); err != nil {
		return err
	}

	id, err := uuid.Parse(in)
	if err != nil {
		return err
	}

	row, err := h.Usecase.FindByID(c.Request().Context(), id)
	if err != nil {
		return err
	}

	out, err := ToMessageDTO(row)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, out)
}
