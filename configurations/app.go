package configurations

import (
	"os"
	"strconv"
)

type appConfiguration struct {
	appName            string
	env                string
	port               string
	corsAllowedOrigins string
	otelCollector      string
	enableTracer       bool
	enableMetric       bool
}

func NewAppConfiguration() *appConfiguration {
	enableTracer, _ := strconv.ParseBool(os.Getenv("ENABLE_TRACER"))
	enableMetric, _ := strconv.ParseBool(os.Getenv("ENABLE_METRIC"))

	return &appConfiguration{
		appName:            os.Getenv("APP_NAME"),
		env:                os.Getenv("APP_ENV"),
		port:               os.Getenv("APP_PORT"),
		corsAllowedOrigins: os.Getenv("CORS_ALLOWED_ORIGINS"),
		otelCollector:      os.Getenv("OTEL_COLLECTOR_URL"),
		enableTracer:       enableTracer,
		enableMetric:       enableMetric,
	}
}

//go:generate mockgen -destination=mocks/mock_app.go -source=app.go IAppConfiguration
type IAppConfiguration interface {
	GetAppName() string
	GetEnv() string
	GetPort() string
	GetCorsAllowedOrigins() string
	GetOtelCollector() string
	GetEnableMetric() bool
	GetEnableTracer() bool
}

func (ac *appConfiguration) GetAppName() string {
	if ac.appName == "" {
		return "kc-ewallet"
	}
	return ac.appName
}

func (ac *appConfiguration) GetEnv() string {
	if ac.env == "" {
		return "local"
	}
	return ac.env
}
func (ac *appConfiguration) GetPort() string {
	if ac.port == "" {
		ac.port = "3006"
	}
	return ac.port
}

func (ac *appConfiguration) GetCorsAllowedOrigins() string {
	return ac.corsAllowedOrigins
}

func (ac *appConfiguration) GetOtelCollector() string {
	return ac.otelCollector
}

func (ac *appConfiguration) GetEnableTracer() bool {
	return ac.enableTracer
}

func (ac *appConfiguration) GetEnableMetric() bool {
	return ac.enableMetric
}
