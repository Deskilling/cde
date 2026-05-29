package core

import (
	_ "embed"
	"fmt"
	"os"

	"charm.land/log/v2"
)

//go:embed shell/cde.zsh
var zshInit string

//go:embed shell/cde.fish
var fishInit string

var shells = map[string]string{
	"zsh":  zshInit,
	"fish": fishInit,
}

func InitShell(shell string) {
	init, ok := shells[shell]
	if !ok {
		log.Errorf("unsupported shell: %s", shell)
		os.Exit(1)
	}
	fmt.Print(init)
}
