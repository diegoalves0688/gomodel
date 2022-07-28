package datadog

import (
	"strconv"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"github.com/diegoalves0688/gomodel/pkg/pathgroup"
	"github.com/labstack/echo/v4"
)

// NewEchoMiddleware returns a "Echo/v4" (github.com/labstack/echo/v4) middleware which will trace
// incoming requests. It accepts a list of path patterns used to filter and ignore.
//
// Paths are expected to follow the pattern "/my/path/:param1/to/:param2".
func NewEchoMiddleware(ignorePaths ...string) echo.MiddlewareFunc {
	ignored := pathgroup.New(ignorePaths)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			request := c.Request()

			// check if current path should be skipped from tracing
			if _, ok := ignored.Find(request.URL.Path); ok {
				return next(c)
			}

			resource := request.Method + " " + c.Path()
			opts := []ddtrace.StartSpanOption{
				tracer.ServiceName(cfg.Tracer.Service),
				tracer.ResourceName(resource),
				tracer.SpanType(ext.SpanTypeWeb),
				tracer.Measured(),
			}

			opts = append(opts, DDTagsFromHTTPRequest(request)...)

			if spanctx, err := tracer.Extract(tracer.HTTPHeadersCarrier(request.Header)); err == nil {
				opts = append(opts, tracer.ChildOf(spanctx))
			}
			var finishOpts []tracer.FinishOption
			span, ctx := tracer.StartSpanFromContext(request.Context(), "http.server", opts...)
			defer func() { span.Finish(finishOpts...) }()

			// pass the span through the request context
			c.SetRequest(request.WithContext(ctx))

			err := next(c)
			if err != nil {
				finishOpts = append(finishOpts, tracer.WithError(err))
				// invokes the registered HTTP error handler
				c.Error(err)
			}

			span.SetTag(ext.HTTPCode, strconv.Itoa(c.Response().Status))
			span.SetTag("http.response.size", strconv.Itoa(int(c.Response().Size)))
			return err
		}
	}
}
