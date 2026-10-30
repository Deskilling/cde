package vscode

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
	"darwin": filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "Code", "User", "globalStorage", "storage.json"),
	"linux":  filepath.Join(core.GetConfigDir(), "Code", "User", "globalStorage", "storage.json"),
}

type storage struct {
	WindowsState struct {
		LastActiveWindow struct {
			Folder string `json:"folder"`
		} `json:"lastActiveWindow"`
	} `json:"windowsState"`
}

func (e *VsCode) ExtractWorkspace() (workspace model.Workspace, err error) {
	storagePath, err := core.CheckOverride(e.Name(), storagePaths)
	if err != nil {
		return model.Workspace{}, fmt.Errorf("failed getting storagePath: %w", err)
	}

	if core.ValidCache(name, storagePath) {
		return core.Cache.Workspaces[name].Workspace, nil
	}

	slog.Debug(storagePath)
	content, err := os.ReadFile(storagePath)
	if err != nil {
		return model.Workspace{}, fmt.Errorf("failed reading file: %w", err)
	}

	var storageJson storage
	err = json.Unmarshal(content, &storageJson)
	if err != nil {
		return model.Workspace{}, fmt.Errorf("failed vscode unmarshal: %w", err)
	}

	path, _ := strings.CutPrefix(storageJson.WindowsState.LastActiveWindow.Folder, "file://")
	slog.Debug(path)

	info, _ := os.Stat(storagePath)
	time := info.ModTime()
	slog.Debug("got time", "time", time.Unix())

	ws := model.Workspace{
		Path:      path,
		Timestamp: time.Unix(),
	}

	err = core.WriteCache(name, core.LatestHash, ws)
	if err != nil {
		return model.Workspace{}, fmt.Errorf("failed vscode to cache: %w", err)
	}

	return ws, nil
}
