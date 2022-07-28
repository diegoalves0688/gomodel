package th

import (
	"context"
	"os"
	"path/filepath"

	"github.com/diegoalves0688/gomodel/pkg/pathutil"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dbfixture"
)

type FixtureHandler func(db *bun.DB) error

func Fixture(ctx context.Context, models []interface{}, file string, opts ...dbfixture.FixtureOption) FixtureHandler {
	return func(db *bun.DB) error {
		db.RegisterModel(models...)

		f := dbfixture.New(db, opts...)
		p := filepath.Join(pathutil.RelativePath(), "..", "..", "testdata")

		return f.Load(ctx, os.DirFS(p), file)
	}
}
