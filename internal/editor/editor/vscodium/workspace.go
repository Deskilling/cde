package vscodium

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cde/internal/editor/model"
	"cde/internal/util"

	"charm.land/log/v2"
)

var storagePaths = map[string]string{
	"darwin": filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "VSCodium", "User", "globalStorage", "storage.json"),
	"linux":  filepath.Join(os.Getenv("HOME"), ".config", "VSCodium", "User", "globalStorage", "storage.json"),
}

type storage struct {
	WindowsState struct {
		LastActiveWindow struct {
			Folder string `json:"folder"`
		} `json:"lastActiveWindow"`
	} `json:"windowsState"`
}

func (e *VsCodium) ExtractWorkspace() (workspace model.Workspace, err error) {
	storagePath, err := util.CheckOverride(e.Name(), storagePaths)
	if err != nil {
		return model.Workspace{}, fmt.Errorf("failed getting storagePath: %w", err)
	}

	log.Debug(storagePath)
	content, err := os.ReadFile(storagePath)
	if err != nil {
		return model.Workspace{}, fmt.Errorf("failed reading file: %w", err)
	}

	var storageJson storage
	json.Unmarshal(content, &storageJson)

	path, _ := strings.CutPrefix(storageJson.WindowsState.LastActiveWindow.Folder, "file://")
	log.Debug(path)

	info, _ := os.Stat(storagePath)
	time := info.ModTime()
	log.Debug(time.Unix())

	return model.Workspace{
		Path:      path,
		Timestamp: time.Unix(),
	}, nil
}
