package db

import (
	"context"
	
	"time"

	"github.com/google/uuid"
)

func UpsertScript(ctx context.Context, name, slug, code string) error {
	path := "/scripts/" + slug

	var exists bool
	err := conn.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM scripts WHERE path = $1
		)
	`, path).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		_, err = conn.Exec(ctx, `
			UPDATE scripts SET name = $1, code = $2, enabled = TRUE WHERE path = $3
		`, name, code, path)
	} else {
		_, err = conn.Exec(ctx, `
			INSERT INTO scripts (id, name, path, code, enabled, created_at)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, uuid.New(), name, path, code, true, time.Now())
	}

	return err
}

type ScriptMeta struct {
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
}

func ListScripts(ctx context.Context) ([]ScriptMeta, error) {
	rows, err := conn.Query(ctx, `
		SELECT name, path, enabled, created_at FROM scripts ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []ScriptMeta
	for rows.Next() {
		var s ScriptMeta
		var path string
		if err := rows.Scan(&s.Name, &path, &s.Enabled, &s.CreatedAt); err != nil {
			return nil, err
		}
		s.Slug = path[len("/scripts/"):]
		results = append(results, s)
	}
	return results, nil
}
