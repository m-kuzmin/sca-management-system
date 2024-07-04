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
INSERT INTO targets (id, name, country)
VALUES (gen_random_uuid(), $1, $2)
RETURNING id;

-- name: LinkTargetToMission :exec
INSERT INTO mission_targets (target_id, mission_id)
VALUES ($1, $2);


-- name: GetTargetsByMissionID :many
WITH target_ids AS (
    SELECT target_id
    FROM mission_targets
    WHERE mission_id = $1
)
SELECT *
FROM targets
WHERE id IN (SELECT target_id FROM target_ids);

-- name: CompleteTarget :exec
UPDATE targets
SET complete = true
WHERE id = $1;
