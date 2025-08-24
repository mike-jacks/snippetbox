package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

type Interval string

const (
	IntervalDays    Interval = "days"
	IntervalWeeks   Interval = "weeks"
	IntervalMonths  Interval = "months"
	IntervalYears   Interval = "years"
	IntervalHours   Interval = "hours"
	IntervalMinutes Interval = "minutes"
	IntervalSeconds Interval = "seconds"
)

func (i Interval) IsValid() bool {
	switch i {
	case IntervalDays, IntervalWeeks, IntervalMonths, IntervalYears, IntervalHours, IntervalMinutes, IntervalSeconds:
		return true
	}
	return false
}

func (m *SnippetModel) Insert(title string, content string, expires map[Interval]int) (int, error) {
	// Build the interval expression dynamically
	var intervalParts []string
	for interval, value := range expires {
		if !interval.IsValid() {
			return 0, fmt.Errorf("invalid interval: %s", interval)
		}
		intervalParts = append(intervalParts, fmt.Sprintf("(%d * INTERVAL '1 %s')", value, interval))
	}

	expires_string := "NULL"
	if len(intervalParts) != 0 {
		intervalExpression := strings.Join(intervalParts, " + ")
		expires_string = fmt.Sprintf("(NOW() AT TIME ZONE 'UTC') + %s", intervalExpression)
	}
	stmt := fmt.Sprintf(`INSERT INTO snippets (title, content, expires) VALUES($1, $2, %s) RETURNING id`, expires_string)
	var id int
	if err := m.DB.QueryRow(stmt, title, content).Scan(&id); err != nil {
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
