package editor

import (
	"errors"
	"os"

	"cde/internal/core"
	"cde/internal/editor/editor/vscode"
	"cde/internal/editor/editor/vscodium"
	"cde/internal/editor/editor/zed"
	"cde/internal/editor/model"

	"charm.land/log/v2"
)

type Loaded struct {
	Editor    model.Editor
	Workspace model.Workspace
}

var Registered []Loaded

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
		if err == nil {
			Registered = append(Registered, Loaded{
				Editor:    editor,
				Workspace: ws,
			})
			continue
		} else {
			log.Warn(err)
		}
	}
}

func Latest() (latest model.Workspace, err error) {
	workingDirectory, _ := os.Getwd()

	for _, v := range Registered {
		if v.Workspace.Path == workingDirectory {
			switch core.GetConfig().Behavior.Repeat {
			case "other":
				return

			case "editor":
				// TODO i need to save the workspace somewhere (maybe in XDG_CACHE_HOME or smth)
				log.Warn("editor not implemented currently")

			case "nothing":

			default:
				log.Warn("invalid Behavior.Repeat key using nothing")
			}
		}

		if v.Workspace.Timestamp > latest.Timestamp {
			log.Debugf("cmp %v(%s) > %v(%s)", v.Workspace.Timestamp, v.Workspace.Path, latest.Timestamp, latest.Path)
			latest = v.Workspace
		}
	}

	if latest.Path == "" {
		return model.Workspace{}, errors.New("no active workspace found")
	}

	return latest, nil
}
