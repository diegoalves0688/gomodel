package logger

import (
	"github.com/diegoalves0688/gomodel/pkg/config"
	tracer "github.com/diegoalves0688/gomodel/pkg/tracer/datadog"
	"go.uber.org/zap"
)

func NewZapLogger(c config.Config) (l *zap.Logger, err error) {
	switch {
	case c.Zap.Preset == config.ZapPresetDevelopment:
		l, err = zap.NewDevelopment()
	case c.Zap.Preset == config.ZapPresetExample:
		l = zap.NewExample()
	default:
		l, err = zap.NewProduction()
	}

	if err == nil {
		l = tracer.NewZapDD(c, l)
	}
	return
}
