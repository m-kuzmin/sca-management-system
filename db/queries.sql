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