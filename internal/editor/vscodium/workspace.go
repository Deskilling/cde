package vscodium

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"cde/internal/editor"

	"charm.land/log/v2"
)

// TODO Support Overrides via config
var storagePaths = map[string]string{
	"linux": filepath.Join(os.Getenv("HOME"), ".config", "VSCodium", "User", "globalStorage", "storage.json"),
}

type storage struct {
	WindowsState struct {
		LastActiveWindow struct {
			Folder string `json:"folder"`
		} `json:"lastActiveWindow"`
	} `json:"windowsState"`
}

func (vscodium *VsCodium) ExtractWorkspace() (workspace editor.Workspace, err error) {
	storagePath, ok := storagePaths[runtime.GOOS]
	if !ok {
		return editor.Workspace{}, fmt.Errorf("unsupported os: %s", runtime.GOOS)
	}

	log.Debug(storagePath)
	content, err := os.ReadFile(storagePath)
	if err != nil {
		return editor.Workspace{}, fmt.Errorf("failed reading file: %w", err)
	}

	var storageJson storage
	json.Unmarshal(content, &storageJson)

	path, _ := strings.CutPrefix(storageJson.WindowsState.LastActiveWindow.Folder, "file://")
	log.Debug(path)

	info, _ := os.Stat(storagePath)
	time := info.ModTime()
	log.Debug(time.Unix())

	return editor.Workspace{
		Path:      path,
		Timestamp: time.Unix(),
	}, nil
}
