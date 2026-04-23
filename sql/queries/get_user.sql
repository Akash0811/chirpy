/* plpgsql-language-server:disable */
-- name: GetUser :one
SELECT *
FROM users
WHERE id=$1;