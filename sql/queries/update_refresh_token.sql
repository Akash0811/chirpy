/* plpgsql-language-server:disable */
-- name: UpdateRefreshToken :exec
UPDATE refresh_tokens
SET updated_at = $2
WHERE token=$1;
