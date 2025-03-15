-- name: GetNotes :many
SELECT *
FROM notes;

-- name: GetNote :one
SELECT * FROM notes
WHERE id = ? LIMIT 1;

-- name: CreateNote :one
INSERT INTO notes (is_complete, created_at, updated_at, note) VALUES (?, ?, ?, ?)
RETURNING *;

-- name: UpdateNote :exec
UPDATE notes 
SET is_complete = ?, updated_at = ?, note = ?
WHERE id = ?;

-- name: DeleteNote :exec
DELETE FROM notes
WHERE id = ?;
