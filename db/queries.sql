-- name: CreateCat :exec
INSERT INTO
    cats (id, name, years_of_experience, breed, salary)
VALUES
    ($1, $2, $3, $4, $5);

-- name: GetCatByID :one
SELECT
    id, name, years_of_experience, breed, salary
FROM
    cats
WHERE
    id = $1;

