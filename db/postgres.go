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

// CreateCat inserts a cat.
func (p *Postgres) CreateCat(ctx context.Context, name string, years_of_experience uint16, breed string, salary uint) (uuid.UUID, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create a uuid: %w", err)
	}

	if err = p.queries.CreateCat(ctx, sqlc.CreateCatParams{
		ID:                id,
		Name:              name,
		YearsOfExperience: int16(years_of_experience),
		Breed:             breed,
		Salary:            int32(salary),
	}); err != nil {
		return uuid.Nil, fmt.Errorf("while creating a cat: %w", err)
	}
	return id, nil
}

// GetCat returns a cat object
func (p *Postgres) GetCatByID(ctx context.Context, id uuid.UUID) (Cat, error) {
	cat, err := p.queries.GetCatByID(ctx, id)
	if err != nil {
		return Cat{}, fmt.Errorf("failed to get a cat with id %s: %w", id, err)
	}

	return Cat{
		Name:              cat.Name,
		Breed:             cat.Breed,
		YearsOfExperience: uint32(cat.YearsOfExperience),
		Salary:            uint64(cat.Salary),
		ID:                cat.ID,
	}, nil
}

func (p *Postgres) GetCatsPaginated(ctx context.Context, amountPerPage, pageNumber uint32) ([]Cat, error) {
	cats, err := p.queries.ListCatsPaginated(ctx, sqlc.ListCatsPaginatedParams{
		Limit:   int32(amountPerPage),
		Column2: pageNumber,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list cats on page %d, with page size %d", pageNumber, amountPerPage)
	}

	result := make([]Cat, len(cats))
	for i, cat := range cats {
		result[i] = Cat{
			Name:              cat.Name,
			Breed:             cat.Breed,
			YearsOfExperience: uint32(cat.YearsOfExperience),
			Salary:            uint64(cat.Salary),
			ID:                cat.ID,
		}
	}

	return result, nil
}

// DeleteCat deletes a cat. Returns nil on successful deletion
func (p *Postgres) DeleteCatByID(ctx context.Context, id uuid.UUID) error {
	if err := p.queries.DeteleCatByID(ctx, id); err != nil {
		return fmt.Errorf("failed to delete a cat with id %s: %w", id, err)
	}

	return nil
}
