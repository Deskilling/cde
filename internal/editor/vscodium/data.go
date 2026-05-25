package vscodium

import "cde/internal/editor"

const name string = "vscodium"

type VsCodium struct{}

func init() {
	editor.Register(&VsCodium{})
}

func (vscodium *VsCodium) Name() string {
	return name
}
