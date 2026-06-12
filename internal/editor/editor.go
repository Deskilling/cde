package editor

import (
	"errors"
	"os"
	"sync"

	"cde/internal/core"
	"cde/internal/editor/editor/vscodium"
	"cde/internal/editor/editor/zed"
	"cde/internal/editor/model"

	"charm.land/log/v2"
)

var Registered []model.Editor

func Load() {
	editors := []model.Editor{
		vscodium.New(),
		zed.New(),
	}

	for _, editor := range editors {
		cfg, ok := core.GetConfig().Editors[editor.Name()]
		if ok && cfg.Disabled {
			continue
		}
		Registered = append(Registered, editor)
	}
}

func Latest() (latest model.Workspace, err error) {
	var mu sync.Mutex
	var wg sync.WaitGroup

	workingDirectory, _ := os.Getwd()

	for _, v := range Registered {
		wg.Add(1)
		go func(editor model.Editor) {
			defer wg.Done()
			w, err := editor.ExtractWorkspace()
			if err != nil {
				return
			}

			if w.Path == workingDirectory {
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

			mu.Lock()
			if w.Timestamp > latest.Timestamp {
				log.Debugf("cmp %v(%s) > %v(%s)", w.Timestamp, w.Path, latest.Timestamp, latest.Path)
				latest = w
			}
			mu.Unlock()
		}(v)
	}
	wg.Wait()

	if latest.Path == "" {
		return model.Workspace{}, errors.New("no active workspace found")
	}

	return latest, nil
}
