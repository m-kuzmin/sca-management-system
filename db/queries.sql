-- name: CreateCat :exec
INSERT INTO
    cats (id, name, years_of_experience, breed, salary)
VALUES
    ($1, $2, $3, $4, $5);

-- name: GetCatByID :one
SELECT *
FROM cats
WHERE id = $1;

-- name: DeteleCatByID :exec
DELETE FROM cats
WHERE id = $1;

-- name: ListCatsPaginated :many
SELECT *
FROM cats
LIMIT $1 OFFSET ($2 - 1) * $1;

-- name: UpdateCatSalaryByID :exec
UPDATE cats
SET salary = $2
WHERE id = $1;

-- name: CreateMission :one
INSERT INTO missions (id)
VALUES (gen_random_uuid())
RETURNING id;

-- name: GetMissionByID :one
SELECT *
FROM missions
WHERE id = $1;

-- name: ListMissions :many
SELECT *
FROM missions
LIMIT $1 OFFSET ($2 - 1) * $1;

-- name: CompleteMission :exec
UPDATE missions
SET complete = true
WHERE id = $1;

-- name: CreateTarget :one
INSERT INTO targets (id, name, country, mission_id)
VALUES (gen_random_uuid(), $1, $2, $3)
RETURNING id;

-- name: GetTargetsByMissionID :many
SELECT *
FROM targets
WHERE mission_id = $1;

-- name: CompleteTarget :exec
UPDATE targets
SET complete = true
WHERE id = $1;

-- name: UpdateTargetNotes :exec
UPDATE targets
SET notes = $2
WHERE id = $1;

-- name: GetTargetCompleteStatus :one
SELECT targets.complete OR missions.complete
FROM targets
JOIN missions ON missions.id = targets.mission_id
WHERE targets.id = $1;

-- name: CountMissionTargets :one
SELECT COUNT(*)
FROM targets
WHERE mission_id = $1;

-- name: AssignCatToMission :exec
UPDATE missions
SET assigned_cat = $2
WHERE id = $1;

-- name: DeleteTarget :exec
DELETE FROM targets
WHERE id = $1;

-- name: DeleteMission :exec
DELETE FROM missions
WHERE id = $1 ;

-- name: IsAssignedMission :one
SELECT assigned_cat IS NOT NULL
FROM missions
WHERE id = $1;

-- name: MissionIsCompleted :one
SELECT complete
FROM missions
WHERE id = $1;