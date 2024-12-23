-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (
  token, 
  created_at, 
  updated_at, 
  user_id, 
  expires_at
)
VALUES ($1, NOW(), NOW(), $2, NOW() + interval '60 days')
RETURNING *;

-- name: GetRefreshToken :one
SELECT 
  *,
  expires_at < NOW() as is_expired
FROM refresh_tokens 
WHERE token = $1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1;
