package main

import (
	"io"
	"kc-ewallet/configurations"
	"kc-ewallet/domains/repository/postgres"
	"kc-ewallet/domains/usecase/transaction"
	"kc-ewallet/domains/usecase/user"
	"kc-ewallet/internals/database"
	"kc-ewallet/internals/helpers/logging"
	rate_limit "kc-ewallet/internals/helpers/rate_limiter"
	"kc-ewallet/internals/helpers/server"
	"kc-ewallet/migrations"
	"kc-ewallet/protocols/http/controller"
	"kc-ewallet/protocols/http/routes"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/zap"
)

func main() {
	// Load environment variables from .env file
	godotenv.Load()

	// Setup logger
	setupLogger()

	// Initialize configurations
	appConfiguration := configurations.NewAppConfiguration()
	databaseConfiguration := configurations.NewDatabaseWriter()
	jwtConfiguration := configurations.NewJWTConfiguration()
	redisConfiguration := configurations.NewRedisConfiguration()

	// Initialize helpers
	// _ := jwt.NewJWTHelper(jwtConfiguration)

	// Set OpenTelemetry propagator to W3C TraceContext for proper traceparent extraction
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	// Initialize tracer
	// traceExp, err := tracer.NewExporter(context.Background(), appConfiguration.GetEnableTracer(), appConfiguration.GetOtelCollector())
	// if err != nil {
	// 	log.Fatal("create trace exporter failed: %w", err)
	// }
	// tp, err := tracer.NewTraceProvider(traceExp, appConfiguration.GetAppName())
	// if err != nil {
	// 	log.Fatal("create trace provider failed: %w", err)
	// }
	// defer func() {
	// 	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// 	defer cancel()

	// 	_ = tp.Shutdown(ctxWithTimeout)
	// }()
	// otel.SetTracerProvider(tp)
	// _ = tp.Tracer(appConfiguration.GetAppName())

	// Initialize metric
	// metricExp, err := metric.NewExporter(context.Background(), appConfiguration.GetEnableMetric(), appConfiguration.GetOtelCollector())
	// if err != nil {
	// 	log.Fatal("create metric exporter failed: %w", err)
	// }

	// mp, err := metric.NewMeterProvider(metricExp, appConfiguration.GetAppName())
	// if err != nil {
	// 	log.Fatal("create metric provider failed: %w", err)
	// }
	// defer func() {
	// 	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// 	defer cancel()

	// 	_ = mp.Shutdown(ctxWithTimeout)
	// }()

	// otel.SetMeterProvider(mp)
	// meter := mp.Meter(appConfiguration.GetAppName())

	// _, err = metric.NewMetric(meter)
	// if err != nil {
	// 	log.Fatal("create new metrics failed: %w", err)
	// }

	log.Printf("trace available: %v at: %s \n", appConfiguration.GetEnableTracer(), appConfiguration.GetOtelCollector())
	log.Printf("metric available: %v at: %s \n", appConfiguration.GetEnableMetric(), appConfiguration.GetOtelCollector())

	// Initialize redis
	rate_limit.InitCoreRedis(redisConfiguration) // for rate limiter
	log.Println("Connected to Redis")

	// Initialize databases
	postgresWriter, err := database.NewPostgresWriter(databaseConfiguration)
	if err != nil {
		log.Fatalf("failed to connect to postgres writer: %v", err)
	}

	// Check and run migration
	if err := migrations.MigrateAll(postgresWriter.GetDB(), databaseConfiguration.GetDBName()); err != nil {
		log.Panicf("error on migrating database: %v", err)
		return
	}

	// Initialize repositories
	postgresRepo := postgres.New(postgresWriter.GetDB())

	// Initialize usecases
	userUsecase := user.NewUserUsecase(postgresWriter.GetDB(), postgresRepo, jwtConfiguration, nil)
	transactionUsecase := transaction.NewTransactionUsecase(postgresWriter.GetDB(), postgresRepo, nil)

	// Initialize controllers
	userController := controller.NewUserController(userUsecase)
	transactionController := controller.NewTransactionController(transactionUsecase)

	// Initialize router with middleware
	router := routes.InitRouter(appConfiguration, nil)

	// Register routes
	routes.RegisterUserRoutes(router, jwtConfiguration.GetSigningKey(), userController)
	routes.RegisterTransactionRoutes(router, jwtConfiguration.GetSigningKey(), transactionController)

	// Create and start the server
	port, err := strconv.Atoi(appConfiguration.GetPort())
	if err != nil {
		log.Fatalf("invalid port -> %+v", err)
	}

	restServer := server.InitRESTServer(&server.RESTServerConfiguration{
		Port:   port,
		Domain: getDomain(appConfiguration),
	}, router)

	log.Printf("Starting http server on port %s...", appConfiguration.GetPort())
	restServer.Serve()
}

func getDomain(appConfiguration configurations.IAppConfiguration) string {
	var domain string
	if appConfiguration.GetEnv() == "dev" {
		domain = "0.0.0.0"
	}

	return domain
}

func setupLogger() {
	logging.LogFilePath = "logs/gin.log"
	logger := logging.New()

	zap.ReplaceGlobals(logger)

	gin.DefaultWriter = io.Discard
	gin.DefaultErrorWriter = io.Discard

	// this is for debugging
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(os.Stdout)
}
