package datadog

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	ddtracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type bunQryHook struct {
	service string
}

// NewBunQueryHook returns a "Bun" (github.com/uptrace/bun) Query Hook with tracing capabilities.
func NewBunQueryHook() bun.QueryHook {
	return &bunQryHook{service: cfg.Tracer.Service}
}

func (h *bunQryHook) BeforeQuery(ctx context.Context, qe *bun.QueryEvent) context.Context {
	const (
		bunOperation = "bun.stmt"
		bunQryLimit  = 10000
	)
	query := qe.QueryTemplate

	if len(query) > bunQryLimit {
		query = query[:bunQryLimit]
	}

	opts := []ddtracer.StartSpanOption{
		ddtracer.SpanType(ext.SpanTypeSQL),
		ddtracer.ServiceName(h.service),
		ddtracer.ResourceName(qe.Operation()),
		ddtracer.Tag(ext.DBStatement, query),
		ddtracer.Tag(ext.DBType, dbSystem(qe.DB)),
		ddtracer.Measured(),
	}

	_, ctx = ddtracer.StartSpanFromContext(ctx, bunOperation, opts...)
	return ctx
}

func (h *bunQryHook) AfterQuery(ctx context.Context, evt *bun.QueryEvent) {
	if span, ok := ddtracer.SpanFromContext(ctx); ok {
		if evt.Err != nil {
			switch evt.Err {
			case nil, sql.ErrNoRows, sql.ErrTxDone:
				// ignore
			default:
				span.SetTag(ext.Error, evt.Err)
			}
		} else if evt.Result != nil {
			if n, _ := evt.Result.RowsAffected(); n > 0 {
				span.SetTag("db.rows_affected", n)
			}
		}
		span.Finish()
	}
}

func dbSystem(db *bun.DB) string {
	switch db.Dialect().Name() {
	case dialect.PG:
		return "postgresql"
	case dialect.MySQL:
		return "mysql"
	case dialect.SQLite:
		return "sqlite"
	case dialect.MSSQL:
		return "mssql"
	default:
		return "unknown"
	}
}
