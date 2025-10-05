package routes

import (
	"kc-ewallet/configurations"
	"kc-ewallet/internals/errors"
	"kc-ewallet/protocols/http/middleware"
	"kc-ewallet/protocols/http/response"
	"net/http"

	"github.com/gin-gonic/gin"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

const Healthz string = "/healthz"
const Readyz string = "/readyz"

// InitRouter initializes the Gin router with middleware
func InitRouter(appConfig configurations.IAppConfiguration, tracer trace.Tracer) *gin.Engine {
	if appConfig.GetEnv() == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(middleware.SlogMiddleware()) // Use our custom slog middleware
	router.Use(middleware.RequestIDInjector())
	router.Use(middleware.PanicRecoveryHandler())
	router.Use(middleware.API())
	router.Use(middleware.CORS(appConfig.GetCorsAllowedOrigins()))
	router.Use(middleware.HandleError())
	router.Use(middleware.SecurityCheck())
	router.Use(middleware.StructuredLogFormatter())
	router.Use(middleware.TraceMiddleware(tracer))
	router.Use(middleware.MetricMiddleware(otel.Meter(appConfig.GetAppName())))
	router.Use(gin.Recovery())

	router.MaxMultipartMemory = 8 << 20 // 8 MiB memory for multipart forms

	router.GET(Healthz, func(c *gin.Context) {
		c.String(http.StatusOK, "healthy")
	})

	router.GET(Readyz, func(c *gin.Context) {
		c.String(http.StatusOK, "ready")
	})

	router.NoRoute(func(c *gin.Context) {
		response.RespondError(c, errors.NotFound.New("Endpoint not found."))
	})

	return router
}

// RegisterHealthRoutes registers health check routes
func RegisterHealthRoutes(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})
}
