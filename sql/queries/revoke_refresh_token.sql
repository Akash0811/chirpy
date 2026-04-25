/* plpgsql-language-server:disable */
-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET updated_at = $2,
revoked_at = $2
WHERE token=$1;
