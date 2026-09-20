package vscodium

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"cde/internal/core"
	"cde/internal/editor/model"
)

var storagePaths = map[string]string{
	"darwin": filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "VSCodium", "User", "globalStorage", "storage.json"),
	"linux":  filepath.Join(core.GetConfigDir(), "VSCodium", "User", "globalStorage", "storage.json"),
}

type storage struct {
	WindowsState struct {
		LastActiveWindow struct {
			Folder string `json:"folder"`
		} `json:"lastActiveWindow"`
	} `json:"windowsState"`
}

func (e *VsCodium) ExtractWorkspace() (workspace model.Workspace, err error) {
	storagePath, err := core.CheckOverride(e.Name(), storagePaths)
	if err != nil {
		return model.Workspace{}, fmt.Errorf("failed getting storagePath: %w", err)
	}

	slog.Debug(storagePath)
	content, err := os.ReadFile(storagePath)
	if err != nil {
		return model.Workspace{}, fmt.Errorf("failed reading file: %w", err)
	}

	var storageJson storage
	json.Unmarshal(content, &storageJson)

	path, _ := strings.CutPrefix(storageJson.WindowsState.LastActiveWindow.Folder, "file://")
	if path == "" {
		return model.Workspace{}, fmt.Errorf("invalid path returned",)
	}

	slog.Debug(path)

	info, _ := os.Stat(storagePath)
	time := info.ModTime()
	slog.Debug("got time", "time", time.Unix())

	return model.Workspace{
		Path:      path,
		Timestamp: time.Unix(),
	}, nil
}
