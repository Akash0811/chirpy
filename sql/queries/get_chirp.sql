/* plpgsql-language-server:disable */
-- name: GetChirp :one
SELECT *
FROM chirps
WHERE id=$1;