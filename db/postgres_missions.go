package db

import (
	"context"
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
			return uuid.Nil, fmt.Errorf("failed to create a mission: %w", err)
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
