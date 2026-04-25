/* plpgsql-language-server:disable */
-- name: GetRefreshToken :one
SELECT *
FROM refresh_tokens
WHERE token=$1;