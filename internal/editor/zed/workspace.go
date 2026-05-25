package zed

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"cde/internal/editor"

	"charm.land/log/v2"
	_ "modernc.org/sqlite"
)

// TODO Support Overrides via config
var dbPaths = map[string]string{
	"darwin": filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "Zed", "db", "0-stable", "db.sqlite"),
	"linux":  filepath.Join(os.Getenv("HOME"), ".local", "share", "zed", "db", "0-stable", "db.sqlite"),
}

func (zed *Zed) ExtractWorkspace() (workspace editor.Workspace, err error) {
	dbPath, ok := dbPaths[runtime.GOOS]
	if !ok {
		return editor.Workspace{}, fmt.Errorf("unsupported os: %s", runtime.GOOS)
	}

	dataSource := fmt.Sprintf("file:%s?mode=ro&_journal=wal", dbPath)

	log.Debugf("open with sqlite %s", dataSource)

	db, err := sql.Open("sqlite", dataSource)
	if err != nil {
		return editor.Workspace{}, fmt.Errorf("open of db failed: %w", err)
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
		return editor.Workspace{}, errors.New("workspace not found")
	}
	if err != nil {
		return editor.Workspace{}, fmt.Errorf("query failed: %w", err)
	}

	timestamp, err := time.Parse("2006-01-02 15:04:05", timestampStr)
	if err != nil {
		return editor.Workspace{}, fmt.Errorf("failed to parse timestamp: %w", err)
	}
	log.Debug(timestamp.Unix())

	return editor.Workspace{
		Path:      path,
		Timestamp: timestamp.Unix(),
	}, nil
}
