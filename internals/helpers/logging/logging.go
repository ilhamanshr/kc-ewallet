package logging

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var (
	LogFileMaxSize int    = 500 // in megabytes
	LogFilePath    string = ""  // path to log file, set as early as possible
)

const (
	ServiceNameKey string = "service_name"

	ResourceKey     string = "resource"
	TimeKey         string = "time"
	LevelKey        string = "severity"
	MessageKey      string = "message"
	StacktraceKey   string = "stacktrace"
	ErrorsKey       string = "errors"
	RequestIDLogKey string = "request_id"
)

type RequestIDKey struct{}

func NewCustomEncoderConfig() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = TimeKey
	encoderConfig.LevelKey = LevelKey
	encoderConfig.MessageKey = MessageKey
	encoderConfig.StacktraceKey = StacktraceKey
	encoderConfig.EncodeTime = func(t time.Time, pae zapcore.PrimitiveArrayEncoder) {
		pae.AppendString(t.Format(time.RFC3339))
	}
	return zapcore.NewJSONEncoder(encoderConfig)
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

	requestID, _ := c.Value(RequestIDKey{}).(string)
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

func New() *zap.Logger {
	var logOutput zapcore.WriteSyncer = os.Stdout

	if LogFilePath != "" {
		logFile := zapcore.AddSync(&lumberjack.Logger{
			Filename: LogFilePath,
			MaxSize:  LogFileMaxSize,
		})
		logOutput = zapcore.NewMultiWriteSyncer(os.Stdout, logFile)
	}

	core := zapcore.NewCore(
		NewCustomEncoderConfig(),
		logOutput,
		zapcore.InfoLevel,
	)

	// logger configured to give timestamp, caller, request id, and stacktrace (on error)
	host, err := os.Hostname()
	if err != nil {
		log.Printf("failed to get host name : %+v", err)
	}

	logger := zap.New(core, zap.Fields(
		zap.String(ResourceKey, host),
	), zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))

	return logger
}
