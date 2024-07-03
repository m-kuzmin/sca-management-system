package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

// Querier can do queries to the entire DB
type Querier interface {
	CatQuerier
}

// CatQuerier can only query the cat info
type CatQuerier interface {
	CreateCat(ctx context.Context, name string, years_of_experience uint16, breed string, salary uint) (uuid.UUID, error)
	GetCatByID(ctx context.Context, id uuid.UUID) (Cat, error)
}

type Cat struct {
	Name              string
	Breed             string
	YearsOfExperience uint32
	Salary            uint64
	ID                uuid.UUID
}

func ConnectToDBWithRetry(driver, address string, retries uint, interval time.Duration) (*sql.DB, error) {
	conn, err := sql.Open(driver, address)
	if err != nil {
		return nil, fmt.Errorf("failed to create sql database connection to \"%q\" with driver \"%q\": %w", address,
			driver, err)
	}

	for i := uint(0); i < retries; i++ {
		err = conn.Ping()
		if err != nil {
			if i > 0 {
				log.Printf("Retry %d: Pinging PostgreSQL after %s because it has not started yet", i, interval.String())
			}

			time.Sleep(interval)

			continue
		}

		return conn, nil
	}

	return nil, fmt.Errorf("failed to ping database %q after %d attempts: %w", address, retries, err)
}
