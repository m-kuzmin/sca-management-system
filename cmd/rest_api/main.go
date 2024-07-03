//nolint:wsl // main() looks better this way
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"

	"github.com/m-kuzmin/sca-management-system/api"
	"github.com/m-kuzmin/sca-management-system/db"
)

func main() {
	config, err := NewConfigFromFile("/etc/server/config.toml")
	if err != nil {
		log.Fatalf("Failed to load config file: %s\n", err)
	}

	postgres := MustSetupPostgres(config.Database)
	defer func() {
		err := postgres.Close()
		if err != nil {
			log.Printf("Error closing connection to postgres: %s", err)
		}
	}()
	log.Println("Connected to postgres")

	gin.SetMode(config.Server.Gin.Mode)
	server := api.NewServer(postgres)
	router := api.NewGinRouter(server)
	log.Println("Setup Gin router")

	httpServer := startServer(router, config.Server)
	defer func() {
		err := httpServer.Shutdown(context.Background())
		if err != nil {
			log.Printf("Error shutting down the server: %s", err)
		}
	}()
	log.Println("Started the HTTP server")

	waitForCtrlC()
	log.Println("Shutting down the server")
}

func startServer(handler http.Handler, c ServerConfig) *http.Server {
	server := &http.Server{
		Addr:        fmt.Sprintf(":%d", c.Port),
		Handler:     handler,
		ReadTimeout: time.Duration(c.ReadTimeoutMs) * time.Millisecond,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("HTTP server error: %s", err)
		}

		log.Printf("HTTP server shutdown")
	}()

	return server
}

func MustSetupPostgres(c DatabaseConfig) *db.Postgres {
	conn, err := db.ConnectToDBWithRetry(c.Driver, c.Address, c.PingRetries, time.Duration(c.PingIntervalSecs)*time.Second)
	if err != nil {
		log.Fatalf("failed to connect to PostgreSQL: %s", err)
	}

	postgres := db.NewPostgres(conn)

	if err = postgres.Migrate(c.MigrationsSource, c.DBName); err != nil {
		log.Fatalf("failed to migrate PostgreSQL to latest version: %s", err)
	}

	return postgres
}

func waitForCtrlC() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	<-interrupt
}
