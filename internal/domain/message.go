package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Message struct {
	bun.BaseModel `bun:"table:message,alias:m"`
	ID            uuid.UUID `bun:"id,pk,nullzero"`
	Receiver      string
	Sender        string
	Content       string
	CreatedAt     time.Time `bun:",nullzero,default:current_timestamp"`
	UpdatedAt     time.Time `bun:",nullzero,default:current_timestamp"`
}
