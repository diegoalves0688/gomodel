package th

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/diegoalves0688/gomodel/pkg/config"
	"github.com/diegoalves0688/gomodel/pkg/db"
	"github.com/diegoalves0688/gomodel/pkg/logger"
	"github.com/diegoalves0688/gomodel/pkg/migrations"
	"github.com/diegoalves0688/gomodel/pkg/pathutil"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"go.uber.org/fx"
)

type IntegrationSuite struct {
	suite.Suite
	Ctx               context.Context
	postgresContainer testcontainers.Container
	postgresDSN       string
}

func (s *IntegrationSuite) loadConfig() (c config.Config, err error) {
	c, err = config.LoadConfig()
	if err != nil {
		return
	}
	c.DB.DSN = s.postgresDSN
	return
}

func (s *IntegrationSuite) Run(funcs ...interface{}) {
	os.Setenv("PROFILE", "test")

	rpath := filepath.Join(pathutil.RelativePath(), "..", "..")

	config.InitProfileConfig(rpath, "yaml")

	app := fx.New(
		fx.Provide(
			s.loadConfig,
			logger.NewZapLogger,
			db.ConnectPostgres,
		),
		fx.Invoke(migrations.RunPostgresMigrations(
			fmt.Sprintf("file://%v", filepath.ToSlash(filepath.Join(rpath, "migrations"))),
		)),

		fx.Invoke(funcs...),
	)

	if err := app.Start(s.Ctx); err != nil {
		log.Fatal(err)
	}
}

func (s *IntegrationSuite) SetupSuite() {
	s.Ctx = context.Background()

	dbDSN, pgContainer, err := NewPostgresContainer(s.Ctx)
	require.NoError(s.T(), err)

	s.postgresContainer = pgContainer
	s.postgresDSN = dbDSN
}

func (s *IntegrationSuite) AfterTest(_, _ string) {
	err := s.postgresContainer.Terminate(s.Ctx)
	require.NoError(s.T(), err)
}
