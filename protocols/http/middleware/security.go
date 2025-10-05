package middleware

import (
	"github.com/gin-gonic/gin"
)

func SecurityCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		cspParam := "default-src 'self'; script-src 'self'; style-src 'self' https://fonts.googleapis.com; img-src 'self' data: ; font-src 'self' https://fonts.gstatic.com;"

		c.Header("Content-Security-Policy", cspParam)   // Chrome 25+, Firefox 23+ and Safari 7+
		c.Header("X-Content-Security-Policy", cspParam) // Firefox 4.0+ and Internet Explorer 10+
		c.Header("X-WebKit-CSP", cspParam)              // Chrome 14+ and Safari 6+

		c.Header("X-Frame-Options", "SAMEORIGIN")

		// Set the X-Content-Type-Options header to "nosniff"
		c.Header("X-Content-Type-Options", "nosniff")

		// Set the Permissions-Policy header
		// only allow geolocation & camera
		c.Header("Permissions-Policy", "geolocation=(), camera=()")

		// Suppress the "Server" header by not calling the parent's ServeHTTP method
		c.Header("Server", "")

		// Set the Strict-Transport-Security header to enforce HTTPS for 30 days (in seconds)
		c.Header("Strict-Transport-Security", "max-age=2592000")
	}
}
