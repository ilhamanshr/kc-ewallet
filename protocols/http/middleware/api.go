package middleware

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// API middleware for API
func API() gin.HandlerFunc {
	return func(c *gin.Context) {
		sentry.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetExtra("method", c.Request.Method)
			scope.SetExtra("url", c.Request.URL.String())
			scope.SetExtra("user_agent", c.Request.UserAgent())
			scope.SetExtra("content_type", c.ContentType())
			scope.SetExtra("query_params", c.Request.URL.Query())
			scope.SetTag("request_id", uuid.New().String())
			scope.SetExtra("json_response", nil)
			scope.SetUser(sentry.User{ID: ""})
			scope.SetExtra("body_params", nil)
		})

		bodyBs, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBs))

		if bodyBs != nil {
			var bodyParams map[string]interface{}
			_ = json.Unmarshal(bodyBs, &bodyParams)
			sentry.ConfigureScope(func(scope *sentry.Scope) {
				scope.SetExtra("body_params", bodyParams)
			})
		}

		c.Next()
	}
}
