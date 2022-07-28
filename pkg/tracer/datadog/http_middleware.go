package datadog

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"github.com/diegoalves0688/gomodel/pkg/pathgroup"
)

type wrapResponseWriter struct {
	http.ResponseWriter
	statusCode  int
	wroteHeader bool
	bodySize    int
}

func (w *wrapResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
	// Check after in case there's error handling in the wrapped ResponseWriter.
	if w.wroteHeader {
		return
	}
	w.statusCode = code
	w.wroteHeader = true
}

func (w *wrapResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *wrapResponseWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	w.bodySize += len(b)
	return w.ResponseWriter.Write(b)
}

// NewHTTPMiddleware return a new HTTP middleware that enables tracing
// of HTTP requests. It accepts a list of patterns used to filter and group similar
// requests under the same "resource name". Ignoring paths is also possible.
//
// Paths are expected to follow the pattern "/my/path/:param1/to/:param2".
func NewHTTPMiddleware(groupPaths []string, ignorePaths []string) func(h http.Handler) http.Handler {
	groups := pathgroup.New(groupPaths)
	ignored := pathgroup.New(ignorePaths)

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// check if current path should be skipped from tracing
			if _, ok := ignored.Find(r.URL.Path); ok {
				h.ServeHTTP(w, r)
				return
			}

			// after this point, trace all requests
			wrapWriter := &wrapResponseWriter{ResponseWriter: w}
			group, _ := groups.Find(r.URL.Path)
			resource := fmt.Sprintf("%v %v", r.Method, group)
			opts := []ddtrace.StartSpanOption{
				tracer.SpanType(ext.SpanTypeWeb),
				tracer.Tag(ext.ServiceName, cfg.Tracer.Service),
				tracer.Tag(ext.ResourceName, resource),
				tracer.Measured(),
			}

			opts = append(opts, DDTagsFromHTTPRequest(r)...)

			if spanctx, err := tracer.Extract(tracer.HTTPHeadersCarrier(r.Header)); err == nil {
				opts = append(opts, tracer.ChildOf(spanctx))
			}
			span, ctx := tracer.StartSpanFromContext(r.Context(), "http.server", opts...)
			defer span.Finish()

			h.ServeHTTP(wrapWriter, r.WithContext(ctx))
			span.SetTag(ext.HTTPCode, strconv.Itoa(wrapWriter.statusCode))
			span.SetTag("http.response.size", strconv.Itoa(wrapWriter.bodySize))

			// treat 5XX as errors
			if wrapWriter.statusCode/100 == 5 {
				status := fmt.Sprintf("%d: %s", wrapWriter.statusCode, http.StatusText(wrapWriter.statusCode))
				span.SetTag("http.errors", status)
				span.SetTag(ext.Error, errors.New(status))
			}
		})
	}
}
