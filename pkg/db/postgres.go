package db

import (
	"database/sql"

	"github.com/diegoalves0688/gomodel/pkg/config"
	tracer "github.com/diegoalves0688/gomodel/pkg/tracer/datadog"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func ConnectPostgres(c config.Config) *bun.DB {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(c.DB.DSN)))
	db := bun.NewDB(sqldb, pgdialect.New())
	db.AddQueryHook(tracer.NewBunQueryHook())
	return db
}
