package datadog

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"

	"github.com/diegoalves0688/gomodel/pkg/config"
	ddtracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var (
	once sync.Once
	cfg  config.Config
)

// Start starts the tracer with the given set of options. It only starts the tracer on
// the first call, meaning that subsequent calls are valid but become no-op.
func Start(c config.Config) {
	once.Do(func() {
		cfg = c
		opts := []ddtracer.StartOption{}

		if cfg.Tracer.EnableMockAgent {
			log.Printf("Datadog Tracer INFO: Config: %v", cfg.Tracer)
			log.Printf("Datadog Tracer WARN: Sending application traces to Mock Agent")
			mockAddr := startMockAgent()
			opts = append(opts, ddtracer.WithAgentAddr(mockAddr))
		} else {
			addr := fmt.Sprintf("%v:%v", cfg.Tracer.Host, cfg.Tracer.Port)
			opts = append(
				opts,
				ddtracer.WithAgentAddr(addr),
				ddtracer.WithRuntimeMetrics(),
			)
		}
		opts = append(
			opts,
			ddtracer.WithService(cfg.Tracer.Service),
			ddtracer.WithEnv(cfg.Tracer.Env),
			ddtracer.WithServiceVersion(cfg.Tracer.Version),
			ddtracer.WithDebugMode(cfg.Tracer.EnableDebugMode),
		)

		ddtracer.Start(opts...)
	})
}

// Stop stops the started tracer. Subsequent calls are valid but become no-op.
func Stop() {
	ddtracer.Stop()
}

func startMockAgent() string {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	url := httptest.NewServer(h).URL
	return strings.ReplaceAll(url, "http://", "")
}
