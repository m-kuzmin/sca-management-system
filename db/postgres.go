package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	"github.com/google/uuid"
	"github.com/m-kuzmin/sca-management-system/db/sqlc"
)

type Postgres struct {
	conn    *sql.DB
	queries *sqlc.Queries
}

// CreateCat implements Querier.
func (p *Postgres) CreateCat(ctx context.Context, name string, years_of_experience uint16, breed string, salary uint) (uuid.UUID, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create a uuid: %w", err)
	}

	err = p.queries.CreateCat(ctx, sqlc.CreateCatParams{ID: id, Name: name, YearsOfExperience: int16(years_of_experience), Breed: breed, Salary: int32(salary)})
	if err != nil {
		return uuid.Nil, fmt.Errorf("while creating a cat: %w", err)
	}
	return id, nil
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
