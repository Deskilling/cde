package model

type Workspace struct {
	Path      string
	Timestamp int64
}

type Editor interface {
	Name() string
	ExtractWorkspace() (Workspace, error)
}
