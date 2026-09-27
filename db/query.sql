-- name: GetMovie :one
SELECT * FROM movies
WHERE id = $1 LIMIT 1;


-- name: GetMovieByPrefix :one
SELECT * FROM movies
WHERE prefix = $1 LIMIT 1;

-- name: ListMovies :many
SELECT * FROM movies
ORDER BY title;

-- name: CreateMovie :one
INSERT INTO movies (
  title, prefix, description
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: UpdateMovie :exec
UPDATE movies
  set title = $2,
  description = $3, 
  prefix = $4
WHERE id = $1;

-- name: DeleteMovie :exec
DELETE FROM movies
WHERE id = $1;
