package middleware

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	RequestIDHeader string = "X-Request-ID"
	RequestIDLogKey string = "request_id"
)

type requestIDKey struct{}

func RequestIDInjector() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(strings.ToLower(RequestIDHeader))
		if requestID == "" {
			requestID = uuid.NewString()
			c.Request.Header[RequestIDHeader] = []string{requestID}
		}
		c.Request = c.Request.WithContext(
			context.WithValue(c.Request.Context(), requestIDKey{}, requestID),
		)
		c.Header(RequestIDHeader, requestID)
		c.Next()
	}
}

// this one use to find correlation between mutiple logs
// to trace for logs that correlated between each others
// there will be some zero request id due to async communication
// will assign new one if found no request id
func GetRequestIDFromContext(ctx context.Context) string {
	c := ctx
	if gc, ok := ctx.(*gin.Context); ok {
		c = gc.Request.Context()
	}

	requestID, _ := c.Value(requestIDKey{}).(string)
	if requestID != "" {
		return requestID
	}

	return uuid.New().String()
}

func NewFromContext(ctx context.Context) *zap.Logger {
	requestID := GetRequestIDFromContext(ctx)
	return zap.L().WithOptions(zap.Fields(
		zap.String(RequestIDLogKey, requestID),
	))
}
