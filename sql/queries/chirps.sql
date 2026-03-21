-- name: CreateChirp :one
INSERT INTO chirps (
    body, user_id
) VALUES (
    $1, $2
)
RETURNING *;

-- name: ListChirps :many
SELECT * FROM chirps
ORDER BY created_at;

-- name: GetChirp :one
SELECT * FROM chirps
WHERE id = $1 LIMIT 1;

-- name: GetChirpsForUser :many
SELECT * FROM chirps
WHERE user_id = $1
ORDER BY created_at ASC;

-- name: GetChirpsByEmail :many
SELECT chirps.* FROM chirps
JOIN users on chirps.user_id = users.id
WHERE users.email = $1
ORDER BY chirps.created_at DESC;

-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE id = $1 AND user_id = $2;
