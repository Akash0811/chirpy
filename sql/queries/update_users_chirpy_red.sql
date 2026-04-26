/* plpgsql-language-server:disable */
-- name: UpdateUserChrirpyRed :exec
UPDATE users
SET is_chirpy_red = $2
WHERE id=$1;
