package message

import (
	"time"

	"github.com/diegoalves0688/gomodel/internal/domain"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
)

type CreateMessageDTO struct {
	Receiver string `json:"receiver"`
	Sender   string `json:"sender"`
	Content  string `json:"content"`
}

func (c *CreateMessageDTO) ToEntity() (to domain.Message, err error) {
	err = mapstructure.Decode(c, &to)
	return
}

type MessageDTO struct {
	ID        uuid.UUID
	Receiver  string    `json:"receiver"`
	Sender    string    `json:"sender"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func ToMessageDTO(from *domain.Message) (to MessageDTO, err error) {
	err = mapstructure.Decode(from, &to)
	return
}
