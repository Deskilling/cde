package zed

import "cde/internal/editor"

const name string = "zed"

type Zed struct{}

func init() {
	editor.Register(&Zed{})
}

func (z *Zed) Name() string {
	return name
}
