package editor

import (
	"errors"
	"os"

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
	for _, v := range Registered {
		w, err := v.ExtractWorkspace()
		if err != nil {
			continue
		}

		// TODO add config toggle
		workingDirectory, _ := os.Getwd()
		if w.Path == workingDirectory {
			continue
		}

		if w.Timestamp > latest.Timestamp {
			log.Debugf("cmp %v(%s) > %v(%s)", w.Timestamp, w.Path, latest.Timestamp, latest.Path)
			latest = w
		}
	}
	if latest.Path == "" {
		return model.Workspace{}, errors.New("no active workspace found")
	}

	return latest, nil
}
