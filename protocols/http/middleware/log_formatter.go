package middleware

import (
	"fmt"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const Healthz string = "/healthz"
const Readyz string = "/readyz"

func LogFormatter() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		var errorMsg string
		if param.ErrorMessage != "" {
			errorMsg = fmt.Sprintf("Error: %s", param.ErrorMessage)
		}
		return fmt.Sprintf("[%s] \"%s %s %s %d\" [%s] - %s %s\n",
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.ClientIP,
			errorMsg,
		)
	})
}

var ignoredPath = []string{Healthz, Readyz}

func StructuredLogFormatter() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			start     time.Time = time.Now()
			path      string    = c.FullPath()
			query     string    = c.Request.URL.RawQuery
			method    string    = c.Request.Method
			userAgent string    = c.Request.UserAgent()
			clientIP  string    = c.ClientIP()
			status    int       = c.Writer.Status()
			requestID string    = GetRequestIDFromContext(c)
		)

		if slices.Contains(ignoredPath, path) {
			c.Next()
			return
		}

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		fields := []zapcore.Field{
			zap.Int("status", status),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", clientIP),
			zap.String("user-agent", userAgent),
			zap.Duration("latency", latency),
			zap.String("time", end.UTC().Format(time.RFC3339)),
			zap.String("request_id", requestID),
		}

		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				zap.L().Error(e, fields...)
			}
		} else {
			zap.L().Info(path, fields...)
		}
	}
}
