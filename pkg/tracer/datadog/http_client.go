package datadog

import (
	"fmt"
	"net/http"

	"github.com/diegoalves0688/gomodel/pkg/pathgroup"
	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
)

// NewHTTPClient modifies the given client's transport to augment it with tracing and returns it.
func NewHTTPClient(c *http.Client, groupPaths ...string) *http.Client {
	// assemble a new path group representation
	groups := pathgroup.New(groupPaths)

	// wrap the HTTP client
	return httptrace.WrapClient(
		c,
		httptrace.RTWithServiceName(cfg.Tracer.Service),
		httptrace.RTWithAnalytics(false),
		httptrace.RTWithResourceNamer(func(r *http.Request) string {
			group, _ := groups.Find(r.URL.Path)
			return fmt.Sprintf("%v %v", r.Method, group)
		}),
	)
}
