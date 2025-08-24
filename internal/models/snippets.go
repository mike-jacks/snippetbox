package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

func (s Snippet) CreatedUTC() string {
	return s.Created.UTC().Format("2006-01-02 15:04:05 UTC")
}

func (s Snippet) ExpiresUTC() string {
	return s.Expires.UTC().Format("2006-01-02 15:04:05 UTC")
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets (title, content, expires) VALUES($1, $2, (NOW() AT TIME ZONE 'UTC') + ($3 * INTERVAL '1 days')) RETURNING id`
	var id int
	if err := m.DB.QueryRow(stmt, title, content, expires).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (m *SnippetModel) Get(id int) (Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets where expires > NOW() AND id = $1`
	var s Snippet
	if err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}
	return s, nil
}

func (m *SnippetModel) Latest() ([]Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets where expires > NOW() ORDER BY id DESC LIMIT 10`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snippets []Snippet

	for rows.Next() {
		var s Snippet
		if err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires); err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
