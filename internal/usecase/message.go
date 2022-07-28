package usecase

import (
	"context"

	"github.com/diegoalves0688/gomodel/internal/domain"
	"github.com/diegoalves0688/gomodel/internal/infra/repository"
	"github.com/google/uuid"
)

type MessageUseCase interface {
	List(context.Context) ([]domain.Message, error)
	Create(context.Context, *domain.Message) error
	FindByID(context.Context, uuid.UUID) (*domain.Message, error)
}

type MessageUseCaseImpl struct {
	Repository repository.MessageRepository
}

func NewMessageUseCaseImpl(r repository.MessageRepository) MessageUseCase {
	return &MessageUseCaseImpl{r}
}

func (u *MessageUseCaseImpl) List(ctx context.Context) ([]domain.Message, error) {
	return u.Repository.Find(ctx)
}

func (u *MessageUseCaseImpl) Create(ctx context.Context, m *domain.Message) error {
	return u.Repository.Create(ctx, m)
}

func (u *MessageUseCaseImpl) FindByID(ctx context.Context, id uuid.UUID) (*domain.Message, error) {
	return u.Repository.FindByID(ctx, id)
}
