package migrations

import (
	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/uptrace/bun"
	"go.uber.org/zap"

	// register File System source for loading migration files.
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type DBMigrationHandler func(db *bun.DB, logger *zap.Logger) error

func RunPostgresMigrations(sourceURL string) DBMigrationHandler {
	return func(db *bun.DB, logger *zap.Logger) error {
		driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
		if err != nil {
			return err
		}

		m, err := migrate.NewWithDatabaseInstance(sourceURL, "postgres", driver)
		if err != nil {
			return err
		}

		err = m.Up()
		if err != migrate.ErrNoChange {
			return err
		}

		return nil
	}
}
