-- name: CreateApplication :one
INSERT INTO application (id, name, path, icon, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ListApplication :many
SELECT * FROM application
ORDER BY name;

-- name: GetApplication :one
SELECT * FROM application
WHERE id = $1 LIMIT 1;

-- name: DeleteApplication :exec
DELETE FROM application
WHERE id = $1;