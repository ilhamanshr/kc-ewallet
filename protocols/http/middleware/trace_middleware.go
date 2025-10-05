package middleware

import (
	"fmt"
	"kc-ewallet/internals/helpers/operation"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"

	"go.opentelemetry.io/otel/trace"
)

const traceStateKey string = "client_request"

// make sure put this after RequestIDInjector
// This will create a new trace from the request and propagate
// if traceparent is found
func TraceMiddleware(tr trace.Tracer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// if no trace given, go next
		if tr == nil {
			c.Next()
			return
		}

		var (
			ctx  = otel.GetTextMapPropagator().Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))
			span = trace.SpanFromContext(ctx)
		)

		if !span.SpanContext().IsValid() {
			traceState := trace.TraceState{}
			traceState, err := traceState.Insert(traceStateKey, encodeTracestateValue(c.FullPath()))
			// if error happen won't trace
			if err != nil {
				NewFromContext(ctx).Error(fmt.Sprintf("inject tracestate failed %s", err.Error()))
				c.Next()
				return
			}

			ctx = trace.ContextWithSpanContext(ctx, trace.NewSpanContext(trace.SpanContextConfig{
				TraceState: traceState,
			}))
		}

		ctx = operation.SetOperationName(ctx, c.FullPath())
		ctx, span = tr.Start(ctx, c.FullPath(),
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				attribute.String("operation", c.FullPath()),
				attribute.String("method", c.Request.Method),
				attribute.String("path", c.Request.URL.String()),
				attribute.String(RequestIDHeader, GetRequestIDFromContext(ctx)),
			),
		)
		defer span.End()

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

var invalidTracestateValueRe = regexp.MustCompile(`[^\x20-\x2b\x2d-\x3c\x3e-\x7e]`)

func encodeTracestateValue(value string) string {
	return invalidTracestateValueRe.ReplaceAllString(strings.TrimSpace(value), "_")
}
