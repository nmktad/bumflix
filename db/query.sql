-- name: GetFilm :one
SELECT * FROM films
WHERE id = $1 LIMIT 1;

-- name: ListFilms :many
SELECT * FROM films
ORDER BY created_at;

-- name: InsertFilm :one
INSERT INTO films (
  title, slug, source_key
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: UpdateFilm :exec
UPDATE films
  set title = $2,
  year = $3,
  transcoded_prefix = $4
WHERE id = $1;

-- name: DeleteFilm :exec
DELETE FROM films
WHERE id = $1;
