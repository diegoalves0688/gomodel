package repository

import (
	"context"

	"github.com/diegoalves0688/gomodel/internal/domain"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type MessageRepository interface {
	Find(context.Context) ([]domain.Message, error)
	Create(context.Context, *domain.Message) error
	FindByID(context.Context, uuid.UUID) (*domain.Message, error)
}

type MessageRepositoryImpl struct {
	DB *bun.DB
}

func NewMessageRepositoryImpl(db *bun.DB) MessageRepository {
	return &MessageRepositoryImpl{db}
}

func (r *MessageRepositoryImpl) Find(ctx context.Context) ([]domain.Message, error) {
	rows := []domain.Message{}
	err := r.DB.NewSelect().Model(&rows).Scan(ctx)
	return rows, err
}

func (r *MessageRepositoryImpl) Create(ctx context.Context, m *domain.Message) (err error) {
	_, err = r.DB.NewInsert().Model(m).Returning("*").Exec(ctx)
	return
}

func (r *MessageRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*domain.Message, error) {
	row := &domain.Message{}
	err := r.DB.NewSelect().Model(row).Where("id = ?", id.String()).Scan(ctx)
	return row, err
}
