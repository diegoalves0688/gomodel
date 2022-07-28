//go:build integration
// +build integration

package repository

import (
	"testing"

	"github.com/diegoalves0688/gomodel/internal/domain"
	"github.com/diegoalves0688/gomodel/pkg/th"
	"github.com/stretchr/testify/suite"
	"github.com/uptrace/bun"
	"gotest.tools/assert"
)

type MessageHTTPIntegrationTest struct {
	th.IntegrationSuite
}

func TestInit(t *testing.T) {
	suite.Run(t, new(MessageHTTPIntegrationTest))
}

func (s *MessageHTTPIntegrationTest) Test_GetAllMessages() {
	var models []interface{}
	models = append(models, (*domain.Message)(nil))

	s.Run(
		th.Fixture(s.Ctx, models, "message.yml"),
		func(db *bun.DB) {
			// Arrange
			repo := NewMessageRepositoryImpl(db)

			// Act
			rows, _ := repo.Find(s.Ctx)

			// Assert
			assert.Equal(s.T(), len(rows), 1)
			assert.Equal(s.T(), "Paulo", rows[0].Receiver)
			assert.Equal(s.T(), "Maria", rows[0].Sender)
		},
	)
}
