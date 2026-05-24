package zed

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"charm.land/log/v2"
	_ "modernc.org/sqlite"
)

var dbPaths = map[string]string{
	"darwin": filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "Zed", "db", "0-stable", "db.sqlite"),
}

func ZedExtractWorkspacePath() (path string, err error) {
	dbPath, ok := dbPaths[runtime.GOOS]
	if !ok {
		return "", fmt.Errorf("unsupported os: %s", runtime.GOOS)
	}

	dataSource := fmt.Sprintf("file:%s?mode=ro&_journal=wal", dbPath)

	log.Debugf("open with sqlite %s", dataSource)

	db, err := sql.Open("sqlite", dataSource)
	if err != nil {
		return "", fmt.Errorf("open of db failed: %w", err)
	}
	defer db.Close()

	err = db.QueryRow(`
		SELECT paths
		FROM workspaces
		WHERE paths != ''
		ORDER BY timestamp DESC
		LIMIT 1
	`).Scan(&path)
	if err == sql.ErrNoRows {
		return "", errors.New("workspace not found")
	}
	if err != nil {
		return "", fmt.Errorf("query failed: %w", err)
	}

	return path, nil
}
