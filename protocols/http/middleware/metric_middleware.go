package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func MetricMiddleware(meter metric.Meter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// if no meter given, go next
		if meter == nil {
			c.Next()
			return
		}

		var (
			logger = NewFromContext(c.Request.Context())
			start  = time.Now()
		)

		requestMeter, err := meter.Int64Counter("http_request_count",
			metric.WithDescription("number of request received"),
		)
		if err != nil {
			logger.Error(fmt.Sprintf("request count meter error %v", err))
			c.Next()
			return
		}

		responseTimeMeter, err := meter.Int64Histogram("http_request_latency",
			metric.WithDescription("latency of request"),
			metric.WithUnit("ms"),
		)
		if err != nil {
			logger.Error(fmt.Sprintf("latency time meter error %v", err))
			c.Next()
			return
		}

		c.Next()

		requestMeter.Add(c.Request.Context(), 1,
			metric.WithAttributes(attribute.Int("status_code", c.Writer.Status())),
			metric.WithAttributes(attribute.String("path", c.FullPath())),
			metric.WithAttributes(attribute.String("method", c.Request.Method)),
		)

		responseTimeMeter.Record(c.Request.Context(), time.Since(start).Milliseconds(),
			metric.WithAttributes(attribute.Int("status_code", c.Writer.Status())),
			metric.WithAttributes(attribute.String("path", c.FullPath())),
			metric.WithAttributes(attribute.String("method", c.Request.Method)),
		)
	}
}
