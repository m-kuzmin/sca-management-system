package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

// Querier can do queries to the entire DB
type Querier interface {
	CatQuerier
}

// CatQuerier can only query the cat info
type CatQuerier interface {
	CreateCat(context.Context, Cat) error
}

type Cat struct {
	ID int64 // Primary key

	Name              string
	YearsOfExperience uint64
	Breed             string
	Salary            uint64
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
