package editor

import (
	"errors"

	"charm.land/log/v2"
)

type Workspace struct {
	Path      string
	Timestamp int64
}

type Editor interface {
	Name() string
	ExtractWorkspace() (Workspace, error)
}

var Registered []Editor

func Register(editor Editor) {
	Registered = append(Registered, editor)
}

func Latest() (latest Workspace, err error) {
	for _, v := range Registered {
		w, err := v.ExtractWorkspace()
		if err != nil {
			continue
		}

		if w.Timestamp > latest.Timestamp {
			log.Debugf("cmp %v(%s) > %v(%s)", w.Timestamp, w.Path, latest.Timestamp, latest.Path)
			latest = w
		}
	}
	if latest.Path == "" {
		return Workspace{}, errors.New("no active workspace found")
	}

	return latest, nil
}
