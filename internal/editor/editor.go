package editor

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"cde/internal/core"
	"cde/internal/editor/editor/vscode"
	"cde/internal/editor/editor/vscodium"
	"cde/internal/editor/editor/zed"
	"cde/internal/editor/model"
)

type Loaded struct {
	Editor    model.Editor
	Workspace model.Workspace
}

var Registered []Loaded

// in the future this should probbably be
// async but with 3 editors is slower than just scanning
func Load() {
	editors := []model.Editor{
		vscodium.New(),
		vscode.New(),
		zed.New(),
	}

	for _, editor := range editors {
		cfg, ok := core.GetConfig().Editors[editor.Name()]
		if ok && cfg.Disabled {
			continue
		}

		ws, err := editor.ExtractWorkspace()
		ws.Path = filepath.Clean(ws.Path)
		if err == nil && ws.Path != "" {
			Registered = append(Registered, Loaded{
				Editor:    editor,
				Workspace: ws,
			})
			slog.Debug("got", "editor", editor.Name(), "workspace", ws)
		} else {
			slog.Warn("no workspace found", "editor", editor.Name(), "err", err)
		}
	}
}

func Latest() (model.Workspace, error) {
	wd, err := os.Getwd()
	if err != nil {
		return model.Workspace{}, fmt.Errorf("getting working directory: %w", err)
	}
	wd = filepath.Clean(wd)

	var latest Loaded

	for _, editor := range Registered {
		if editor.Workspace.Timestamp >= latest.Workspace.Timestamp {
			latest = editor
		}
	}

	if latest == (Loaded{}) {
		return model.Workspace{}, errors.New("no active workspace found")
	}

	// always cd even if in the same dir
	return latest.Workspace, nil
}
