//nolint:wsl // main() looks better this way
package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"

	"github.com/m-kuzmin/sca-management-system/db"
)

func main() {
	config, err := NewConfigFromFile("/etc/server/config.toml")
	if err != nil {
		log.Fatalf("Failed to load config file: %s", err)
	}

	postgres := MustSetupPostgres(config.Database)
	defer postgres.Close()
	log.Printf("Connected to postgres")

	gin.SetMode(config.Server.Gin.Mode)
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
