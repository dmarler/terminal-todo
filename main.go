package main

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dmarler/terminal-todo/nindex"
	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type appState int

const (
	visual appState = iota
	adding
	deleting
)

type model struct {
	items   *[]nindex.Note
	input   textinput.Model
	ctx     context.Context
	queries *nindex.Queries
	cursor  int
	state   appState
}

func main() {
	ctx := context.Background()
	db, err := sql.Open("sqlite", "notes.sqlite")
	if err != nil {
		panic(err)
	}

	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("sqlite3"); err != nil {
		panic(err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		panic(err)
	}

	queries := nindex.New(db)

	p := tea.NewProgram(initialModel(ctx, queries))
	if _, err := p.Run(); err != nil {
		fmt.Println("There has been an error.")
		os.Exit(1)
	}

}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		switch m.state {
		case adding:
			return m.updateAdd(msg)
		case visual:
			return m.updateVisual(msg)
		case deleting:
			return m.updateDelete(msg)
		}
	}

	return m, cmd
}

func (m model) updateVisual(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "a":
			m.input.Focus()
			m.state = adding
		case "d":
			m.state = deleting
		case "k", "up":
			if m.cursor > 0 {
				m.cursor -= 1
			}
		case "j", "down":
			if m.cursor <= len(*m.items)-1 {
				m.cursor += 1
			}
		}
	}

	return m, cmd
}

func (m model) updateAdd(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.createNote(m.input.Value())
			m.input.SetValue("")
			m.input.Blur()
			m.fetchNotes()
			m.state = visual
		case "escape":
			m.input.SetValue("")
			m.input.Blur()
		}
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) updateDelete(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "d":
			m.deleteNote((*m.items)[m.cursor].ID)
			m.fetchNotes()
			m.state = visual
			return m, cmd
		case "escape":
			m.state = visual
			return m, cmd
		}
		return m, cmd
	}

	return m, cmd
}

func (m model) View() string {
	var b strings.Builder
	b.WriteString("Terminal Todo\n\n")

	for i, note := range *m.items {
		var cursorValue string
		if m.cursor == i {
			cursorValue = ">"
		} else {
			cursorValue = " "
		}
		b.WriteString(fmt.Sprintf("  %s [] - %s\n", cursorValue, note.Note))
	}

	b.WriteString(m.input.View())

	return b.String()
}

func initialModel(ctx context.Context, queries *nindex.Queries) model {
	m := model{
		items:   nil,
		input:   textinput.New(),
		ctx:     ctx,
		queries: queries,
		cursor:  0,
		state:   visual,
	}

	m.fetchNotes()

	return m
}

func (m model) createNote(note string) {
	if _, err := createNote(m.ctx, m.queries, note); err != nil {
		panic(err)
	}
}

func (m model) deleteNote(id int64) {
	if err := deleteNote(m.ctx, m.queries, id); err != nil {
		panic(err)
	}
}

func (m *model) fetchNotes() {
	notes, err := m.queries.GetNotes(m.ctx)
	if err != nil {
		panic(err)
	}
	m.items = &notes
}

func createNote(ctx context.Context, queries *nindex.Queries, note string) (*nindex.Note, error) {
	t := time.Now()
	nt := sql.NullTime{Time: t, Valid: true}
	noteParams := nindex.CreateNoteParams{
		CreatedAt:  nt,
		UpdatedAt:  nt,
		IsComplete: false,
		Note:       note,
	}
	n, err := queries.CreateNote(ctx, noteParams)
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func updateNote(ctx context.Context, queries *nindex.Queries, id int, isComplete bool, note string) error {
	t := sql.NullTime{Time: time.Now(), Valid: true}

	noteParams := nindex.UpdateNoteParams{
		IsComplete: isComplete,
		UpdatedAt:  t,
		Note:       note,
		ID:         int64(id),
	}

	err := queries.UpdateNote(ctx, noteParams)
	if err != nil {
		return err
	}
	return nil
}

func deleteNote(ctx context.Context, queries *nindex.Queries, id int64) error {
	if err := queries.DeleteNote(ctx, int64(id)); err != nil {
		panic(err)
	}
	return nil
}
