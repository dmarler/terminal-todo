// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: notes.sql

package nindex

import (
	"context"
	"database/sql"
)

const createNote = `-- name: CreateNote :one
INSERT INTO notes (is_complete, created_at, updated_at, note) VALUES (?, ?, ?, ?)
RETURNING id, is_complete, created_at, updated_at, note
`

type CreateNoteParams struct {
	IsComplete bool
	CreatedAt  sql.NullTime
	UpdatedAt  sql.NullTime
	Note       string
}

func (q *Queries) CreateNote(ctx context.Context, arg CreateNoteParams) (Note, error) {
	row := q.db.QueryRowContext(ctx, createNote,
		arg.IsComplete,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Note,
	)
	var i Note
	err := row.Scan(
		&i.ID,
		&i.IsComplete,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Note,
	)
	return i, err
}

const deleteNote = `-- name: DeleteNote :exec
DELETE FROM notes
WHERE id = ?
`

func (q *Queries) DeleteNote(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteNote, id)
	return err
}

const getNote = `-- name: GetNote :one
SELECT id, is_complete, created_at, updated_at, note FROM notes
WHERE id = ? LIMIT 1
`

func (q *Queries) GetNote(ctx context.Context, id int64) (Note, error) {
	row := q.db.QueryRowContext(ctx, getNote, id)
	var i Note
	err := row.Scan(
		&i.ID,
		&i.IsComplete,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Note,
	)
	return i, err
}

const getNotes = `-- name: GetNotes :many
SELECT id, is_complete, created_at, updated_at, note
FROM notes
`

func (q *Queries) GetNotes(ctx context.Context) ([]Note, error) {
	rows, err := q.db.QueryContext(ctx, getNotes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Note
	for rows.Next() {
		var i Note
		if err := rows.Scan(
			&i.ID,
			&i.IsComplete,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Note,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateNote = `-- name: UpdateNote :exec
UPDATE notes 
SET is_complete = ?, updated_at = ?, note = ?
WHERE id = ?
`

type UpdateNoteParams struct {
	IsComplete bool
	UpdatedAt  sql.NullTime
	Note       string
	ID         int64
}

func (q *Queries) UpdateNote(ctx context.Context, arg UpdateNoteParams) error {
	_, err := q.db.ExecContext(ctx, updateNote,
		arg.IsComplete,
		arg.UpdatedAt,
		arg.Note,
		arg.ID,
	)
	return err
}
