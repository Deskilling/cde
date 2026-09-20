package zed

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"cde/internal/core"
	"cde/internal/editor/model"

	_ "modernc.org/sqlite"
)

var dbPaths = map[string]string{
	"darwin": filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "Zed", "db", "0-stable", "db.sqlite"),
	"linux":  filepath.Join(os.Getenv("HOME"), ".local", "share", "zed", "db", "0-stable", "db.sqlite"),
}

func (e *Zed) ExtractWorkspace() (workspace model.Workspace, err error) {
	dbPath, err := core.CheckOverride(e.Name(), dbPaths)
	if err != nil {
		return model.Workspace{}, fmt.Errorf("failed getting dbPath: %w", err)
	}

	dataSource := fmt.Sprintf("file:%s?mode=ro&_journal=wal", dbPath)

	slog.Debug("open with sqlite","path", dataSource)

	db, err := sql.Open("sqlite", dataSource)
	if err != nil {
		return model.Workspace{}, fmt.Errorf("open of db failed: %w", err)
	}
	defer db.Close()

	var path string
	var timestampStr string

	err = db.QueryRow(`
		SELECT paths, timestamp
		FROM workspaces
		WHERE paths != ''
		ORDER BY timestamp DESC
		LIMIT 1
	`).Scan(&path, &timestampStr)
	if err == sql.ErrNoRows {
		return model.Workspace{}, errors.New("workspace not found")
	}
	if err != nil {
		return model.Workspace{}, fmt.Errorf("query failed: %w", err)
	}

	timestamp, err := time.Parse("2006-01-02 15:04:05", timestampStr)
	if err != nil {
		return model.Workspace{}, fmt.Errorf("failed to parse timestamp: %w", err)
	}
	slog.Debug("got time", "time", timestamp.Unix())

	return model.Workspace{
		Path:      path,
		Timestamp: timestamp.Unix(),
	}, nil
}
