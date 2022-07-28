package datadog

import (
	"context"

	"github.com/diegoalves0688/gomodel/pkg/config"
	"go.uber.org/zap"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Start returns a new logger injected with default additional info
// required by Datadog's tracing platform.
func NewZapDD(c config.Config, l *zap.Logger) *zap.Logger {
	return l.With(
		zap.String("service", c.Tracer.Service),
		zap.String("version", c.Tracer.Version),
		// dd is prefix to DataDog (current APM solution)
		zap.String("dd.service", c.Tracer.Service),
		zap.String("dd.env", c.Tracer.Env),
		zap.String("dd.version", c.Tracer.Version),
	)
}

// L returns a new logger with additional tracing info based on the given
// input context.
func L(ctx context.Context, l *zap.Logger) *zap.Logger {
	if spanFromContext, valid := tracer.SpanFromContext(ctx); valid {
		return l.With(
			zap.Uint64("dd.trace_id", spanFromContext.Context().TraceID()),
			zap.Uint64("dd.span_id", spanFromContext.Context().SpanID()),
		)
	}

	return l
}
