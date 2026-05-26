package editor

import (
	"errors"
	"os"
	"sync"

	"cde/internal/editor/editor/vscodium"
	"cde/internal/editor/editor/zed"
	"cde/internal/editor/model"

	"charm.land/log/v2"
)

var Registered []model.Editor

func Load() {
	Registered = append(Registered, vscodium.New())
	Registered = append(Registered, zed.New())
}

func Latest() (latest model.Workspace, err error) {
	var mu sync.Mutex
	var wg sync.WaitGroup

	workingDirectory, _ := os.Getwd()

	for _, v := range Registered {
		wg.Add(1)
		go func(editor model.Editor) {
			defer wg.Done()
			w, err := v.ExtractWorkspace()
			if err != nil {
				return
			}

			// TODO add config toggle
			if w.Path == workingDirectory {
				return
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
