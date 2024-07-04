package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/m-kuzmin/sca-management-system/db/sqlc"
)

func (p *Postgres) CreateMission(ctx context.Context) (uuid.UUID, error) {
	id, err := p.queries.CreateMission(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("while creating a mission: %w", err)
	}

	return id, nil
}

/*
CreateMissionWithTargets creates a mission and adds targets in a transaction. If the targets is nil or [],
calls CreateMission.
*/
func (p *Postgres) CreateMissionWithTargets(ctx context.Context, targets []CreateTargetParams) (uuid.UUID, error) {
	if len(targets) == 0 {
		return p.CreateMission(ctx)
	}

	tx, err := p.conn.Begin()
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	qtx := p.queries.WithTx(tx)

	missionID, err := qtx.CreateMission(ctx)
	if err != nil {
		return uuid.Nil, err
	}

	for _, t := range targets {
		targetID, err := qtx.CreateTarget(ctx, sqlc.CreateTargetParams{
			Name:    t.Name,
			Country: t.Country,
		})
		if err != nil {
			return uuid.Nil, fmt.Errorf("failed to create a target: %w", err)
		}

		err = qtx.LinkTargetToMission(ctx, sqlc.LinkTargetToMissionParams{
			MissionID: missionID,
			TargetID:  targetID,
		})

		if err != nil {
			return uuid.Nil, fmt.Errorf("failed to link the mission and the target: %w", err)
		}
	}

	return missionID, tx.Commit()
}

// GetMissionByID returns a mission with a matching id
func (p *Postgres) GetMissionByID(ctx context.Context, id uuid.UUID) (MissionWithTargets, error) {
	dbMission, err := p.queries.GetMissionByID(ctx, id)
	if err != nil {
		return MissionWithTargets{}, fmt.Errorf("failed to get a mission with id %s: %w", id, err)
	}

	dbTargets, err := p.queries.GetTargetsByMissionID(ctx, dbMission.ID)
	if err != nil {
		return MissionWithTargets{}, fmt.Errorf("failed to get mission %s targets", id)
	}

	targets := make([]Target, len(dbTargets))
	for i, t := range dbTargets {
		targets[i] = Target{
			ID:       t.ID,
			Name:     t.Name,
			Country:  t.Country,
			Notes:    nil,
			Complete: t.Complete,
		}

		if t.Notes.Valid {
			targets[i].Notes = &t.Notes.String
		}
	}

	mission := MissionWithTargets{
		ID:          dbMission.ID,
		Complete:    dbMission.Complete,
		AssignedCat: nil,
		Targets:     targets,
	}

	if dbMission.AssignedCat.Valid {
		mission.AssignedCat = &dbMission.AssignedCat.UUID
	}

	return mission, nil
}

func (p *Postgres) ListMissions(ctx context.Context, pagination PaginationParams) ([]Mission, error) {
	missions, err := p.queries.ListMissions(ctx, sqlc.ListMissionsParams{
		Limit:   int32(pagination.Limit),
		Column2: pagination.PageNumber,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get the mission list: %w", err)
	}

	result := make([]Mission, len(missions))
	for i, m := range missions {
		result[i] = Mission{
			ID:          m.ID,
			AssignedCat: nil,
			Complete:    m.Complete,
		}

		if m.AssignedCat.Valid {
			result[i].AssignedCat = &m.AssignedCat.UUID
		}
	}

	return result, nil
}

func (p *Postgres) CompleteMission(ctx context.Context, id uuid.UUID) error {
	err := p.queries.CompleteMission(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to complete mission with id %s: %w", id, err)
	}

	return nil
}

func (p *Postgres) CountMissionTargets(ctx context.Context, id uuid.UUID) (uint64, error) {
	count, err := p.queries.CountMissionTargets(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("failed to count targets: %w", err)
	}

	return uint64(count), nil
}

func (p *Postgres) AddTargetsToMission(ctx context.Context, missionID uuid.UUID, targets []CreateTargetParams) ([]uuid.UUID, error) {
	tx, err := p.conn.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	qtx := p.queries.WithTx(tx)

	ids := make([]uuid.UUID, len(targets))
	for i, t := range targets {
		targetID, err := qtx.CreateTarget(ctx, sqlc.CreateTargetParams{
			Name:    t.Name,
			Country: t.Country,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create a target: %w", err)
		}

		err = qtx.LinkTargetToMission(ctx, sqlc.LinkTargetToMissionParams{
			MissionID: missionID,
			TargetID:  targetID,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to link the mission and the target: %w", err)
		}

		ids[i] = targetID
	}

	return ids, tx.Commit()
}

func (p *Postgres) GetTargetCompleteStatus(ctx context.Context, id uuid.UUID) (bool, error) {
	complete, err := p.queries.GetTargetCompeleteStatus(ctx, id)
	if err != nil {
		return false, fmt.Errorf("failed to complete target with id %s: %w", id, err)
	}

	return complete, nil
}

func (p *Postgres) CompleteTarget(ctx context.Context, id uuid.UUID) error {
	err := p.queries.CompleteTarget(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to complete target with id %s: %w", id, err)
	}

	return nil
}

func (p *Postgres) UpdateTargetNotes(ctx context.Context, id uuid.UUID, notes string) error {
	err := p.queries.UpdateTargetNotes(ctx, sqlc.UpdateTargetNotesParams{
		ID: id,
		Notes: sql.NullString{
			Valid:  len(notes) == 0,
			String: notes,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to complete target with id %s: %w", id, err)
	}

	return nil
}

func (p *Postgres) AssignCatToMission(ctx context.Context, params AssignCatToMissionParams) error {
	err := p.queries.AssignCatToMission(ctx, sqlc.AssignCatToMissionParams{
		ID:          params.Mission,
		AssignedCat: uuid.NullUUID{Valid: true, UUID: params.Cat},
	})
	if err != nil {
		return fmt.Errorf("failed to assign a cat to mission: %w", err)
	}

	return nil
}
