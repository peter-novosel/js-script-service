package db

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/peter-novosel/js-script-service/internal/logger"
)

type Script struct {
	ID        uuid.UUID
	Name      string
	Path      string
	Code      string
	Enabled   bool
	CreatedAt time.Time
}

func GetScriptByPath(ctx context.Context, path string) (*Script, error) {
	row := conn.QueryRow(ctx, `
		SELECT id, name, path, code, enabled, created_at
		FROM scripts
		WHERE path = $1 AND enabled = TRUE
	`, path)

	var s Script
	err := row.Scan(&s.ID, &s.Name, &s.Path, &s.Code, &s.Enabled, &s.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("script not found")
		}
		return nil, err
	}

	return &s, nil
}

func LogExecution(ctx context.Context, scriptID uuid.UUID, input interface{}, output interface{}, errMsg string) {
	inputJSON, _ := json.Marshal(input)
	outputJSON, _ := json.Marshal(output)

	_, err := conn.Exec(ctx, `
		INSERT INTO execution_logs (id, script_id, input, output, error)
		VALUES ($1, $2, $3, $4, $5)
	`, uuid.New(), scriptID, inputJSON, outputJSON, errMsg)

	if err != nil {
		log := logger.Init()
		log.WithError(err).Error("failed to log script execution")
	}
}
