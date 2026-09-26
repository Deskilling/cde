package godot

import (
	"cde/internal/core"
	"cde/internal/editor/model"

	"bufio"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

var storagePaths = map[string]string{
	"darwin": filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "Godot"),
}

func (e *Godot) ExtractWorkspace() (workspace model.Workspace, err error) {
	storagePath, err := core.CheckOverride(e.Name(), storagePaths)
	if err != nil {
		return model.Workspace{}, fmt.Errorf("failed getting storagePath: %w", err)
	}

	file, err := os.Open(filepath.Join(storagePath, "projects.cfg"))
	if err != nil {
		return model.Workspace{}, fmt.Errorf("failed opening file %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Scan()
	path := scanner.Text()

	err = scanner.Err()
	if err != nil {
		return model.Workspace{}, fmt.Errorf("scanner failed: %w", err)
	}

	path = strings.TrimPrefix(path, "[")
	path = strings.TrimSuffix(path, "]")

	entries, err := os.ReadDir(storagePath)
	if err != nil {
		return model.Workspace{}, fmt.Errorf("failed reading directory: %w", err)
	}

	var matches []os.DirEntry
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "editor_settings-") {
			matches = append(matches, e)
		}
	}

	if len(matches) == 0 {
		return model.Workspace{}, fmt.Errorf("no editor_settings-.* found")
	}
	slog.Debug(filepath.Join(storagePath, matches[0].Name()))
	info, _ := os.Stat(filepath.Join(storagePath, matches[0].Name()))

	time := info.ModTime()
	slog.Debug("got time", "time", time.Unix())

	return model.Workspace{
		Path:      path,
		Timestamp: time.Unix(),
	}, nil
}
