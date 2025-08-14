package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int       `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Created time.Time `json:"created"`
	Expires time.Time `json:"expires"`
}

type SnippetModel struct {
	DB *sql.DB
}

func (model *SnippetModel) Insert(title, content string, expires int) (int, error) {
	queryStatement := `
		INSERT INTO snippets (title, content, created, expires)
		VALUES (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))
	`
	result, err := model.DB.Exec(queryStatement, title, content, expires)
	if err != nil {
		return -1, nil
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, nil
	}
	return int(id), nil
}

func (model *SnippetModel) Get(id int) (*Snippet, error) {
	queryStatement := `
		SELECT * FROM snippets
		WHERE expires > UTC_TIMESTAMP() AND ID = ?
	`

	row := model.DB.QueryRow(queryStatement, id)

	snippet, err := scanRowIntoSnippet(row)
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return snippet, nil
}

func (model *SnippetModel) Latest() ([]*Snippet, error) {
	snippets := make([]*Snippet, 0, 10)
	queryStatement := `
		SELECT * FROM snippets
		WHERE expires > UTC_TIMESTAMP()
		ORDER BY id DESC
		LIMIT 10
	`
	rows, err := model.DB.Query(queryStatement)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		snippet, err := scanRowIntoSnippet(rows)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, snippet)
	}

	// check if an error occured during the iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanRowIntoSnippet(row scanner) (*Snippet, error) {
	snippet := new(Snippet)
	err := row.Scan(
		&snippet.ID,
		&snippet.Title,
		&snippet.Content,
		&snippet.Created,
		&snippet.Expires,
	)
	if err != nil {
		return nil, err
	}
	return snippet, nil
}
