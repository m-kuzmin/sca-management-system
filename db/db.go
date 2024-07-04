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
	MissionQuerier
}

// CatQuerier can only query the cat info
type CatQuerier interface {
	CreateCat(ctx context.Context, name string, years_of_experience uint16, breed string, salary uint) (uuid.UUID, error)
	GetCatByID(context.Context, uuid.UUID) (Cat, error)
	GetCatsPaginated(context.Context, PaginationParams) ([]Cat, error)
	UpdateCatSalaryByID(context.Context, uuid.UUID, uint32) error
	DeleteCatByID(context.Context, uuid.UUID) error
}

type PaginationParams struct {
	PageNumber, Limit uint32
}

type Cat struct {
	Name              string    `json:"name"`
	Breed             string    `json:"breed"`
	YearsOfExperience uint32    `json:"years_of_experience"`
	Salary            uint64    `json:"salary"`
	ID                uuid.UUID `json:"id"`
}

type MissionQuerier interface {
	CreateMission(ctx context.Context) (uuid.UUID, error)
	CreateMissionWithTargets(ctx context.Context, targets []CreateTargetParams) (uuid.UUID, error)
	GetMissionByID(context.Context, uuid.UUID) (MissionWithTargets, error)
	ListMissions(context.Context, PaginationParams) ([]Mission, error)
	CompleteMission(context.Context, uuid.UUID) error
	AddTargetsToMission(ctx context.Context, missionID uuid.UUID, _ []CreateTargetParams) ([]uuid.UUID, error)
	GetTargetCompleteStatus(ctx context.Context, id uuid.UUID) (bool, error)
	CompleteTarget(ctx context.Context, id uuid.UUID) error
	UpdateTargetNotes(ctx context.Context, id uuid.UUID, notes string) error
	CountMissionTargets(context.Context, uuid.UUID) (uint64, error)
	AssignCatToMission(ctx context.Context, params AssignCatToMissionParams) error
	DeleteTarget(ctx context.Context, id uuid.UUID) error
	DeleteMission(ctx context.Context, id uuid.UUID) error
	IsAssignedMission(ctx context.Context, id uuid.UUID) (bool, error)
	MissionIsCompleted(ctx context.Context, id uuid.UUID) (bool, error)
}

type MissionWithTargets struct {
	ID          uuid.UUID  `json:"id"`
	AssignedCat *uuid.UUID `json:"assigned_cat"`
	Complete    bool       `json:"complete"`
	Targets     []Target   `json:"targets"`
}

type Mission struct {
	ID          uuid.UUID  `json:"id"`
	AssignedCat *uuid.UUID `json:"assigned_cat"`
	Complete    bool       `json:"complete"`
}

type AssignCatToMissionParams struct {
	Cat     uuid.UUID
	Mission uuid.UUID
}

type CreateTargetParams struct {
	Name    string `json:"name"`
	Country string `json:"country"`
}

type Target struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Country  string    `json:"country"`
	Notes    *string   `json:"notes"`
	Complete bool      `json:"complete"`
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
