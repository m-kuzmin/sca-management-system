package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	"github.com/m-kuzmin/sca-management-system/db/sqlc"
)

type Postgres struct {
	conn    *sql.DB
	queries *sqlc.Queries
}

func NewPostgres(conn *sql.DB) *Postgres {
	return &Postgres{
		conn:    conn,
		queries: sqlc.New(conn),
	}
}

func (p *Postgres) Migrate(migrationsSource, dbName string) error {
	driver, err := postgres.WithInstance(p.conn, &postgres.Config{
		MigrationsTable: "migrations",
		DatabaseName:    dbName,
	})
	if err != nil {
		return fmt.Errorf("failed to create PostgreSQL database driver: %w", err)
	}

	migrator, err := migrate.NewWithDatabaseInstance(migrationsSource, dbName, driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate client: %w", err)
	}

	if err = migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to migrate the database: %w", err)
	}

	return nil
}

func (p *Postgres) Close() error {
	if err := p.conn.Close(); err != nil {
		return fmt.Errorf("failed to close sql connection handle: %w", err)
	}

	return nil
}
