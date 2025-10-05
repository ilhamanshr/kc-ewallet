package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type RESTServerConfiguration struct {
	Port   int
	Domain string
}

// RESTServer represents http REST server
type RESTServer struct {
	router http.Handler
	config *RESTServerConfiguration
}

// Serve serves the http requests to http server
func (hs *RESTServer) Serve() {
	port := hs.config.Port
	if os.Getenv("SERVER_PORT") != "" {
		parsedPort, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
		if err == nil {
			port = parsedPort
		}
	}

	address := fmt.Sprintf("%s:%d", hs.config.Domain, port)

	srv := &http.Server{
		ReadTimeout:       1 * time.Minute,
		ReadHeaderTimeout: 20 * time.Second,
		Addr:              address,
		Handler:           hs.router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)

	// The server will listen to the SIGINT and SIGTERM
	// SIGINT will listen to CTRL-C.
	// SIGTERM will be caught if kill command executed.
	//
	// See:
	// - https://en.wikipedia.org/wiki/Unix_signal
	// - https://www.quora.com/What-is-the-difference-between-the-SIGINT-and-SIGTERM-signals-in-Linux
	// - http://programmergamer.blogspot.co.id/2013/05/clarification-on-sigint-sigterm-sigkill.html
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

// InitRESTServer initialize REST server & setup routes
func InitRESTServer(config *RESTServerConfiguration, router *gin.Engine) *RESTServer {
	srv := &RESTServer{
		config: config,
		router: router,
	}

	return srv
}
