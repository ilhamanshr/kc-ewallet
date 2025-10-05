package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// SlogMiddleware returns a Gin middleware that logs requests using slog
func SlogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Skip logging for health check endpoints
		if path == "/healthz" || path == "/readyz" || path == "/health" {
			return
		}

		// Calculate latency
		latency := time.Since(start)

		// Log request details
		slog.Info("HTTP Request",
			"status", c.Writer.Status(),
			"method", c.Request.Method,
			"path", path,
			"query", raw,
			"ip", c.ClientIP(),
			"latency", latency,
			"user-agent", c.Request.UserAgent(),
			"errors", c.Errors.String(),
		)
	}
}
