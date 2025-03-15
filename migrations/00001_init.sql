-- +goose Up
-- +goose StatementBegin
CREATE TABLE notes (
  id INTEGER PRIMARY KEY,
  is_complete BOOLEAN not null DEFAULT false,
  created_at DATE,
  updated_at DATE,
  note TEXT not null
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE notes;
-- +goose StatementEnd
