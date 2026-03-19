-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (
    token, user_id, expires_at
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetToken :one
SELECT * FROM refresh_tokens
WHERE token = $1
AND revoked_at IS NULL
AND expires_at > NOW()
LIMIT 1;

-- name: RevokeToken :exec
UPDATE refresh_tokens
    set revoked_at = NOW(),
    updated_at = NOW()
WHERE token = $1;
