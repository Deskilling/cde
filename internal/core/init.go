package core

import (
	_ "embed"
	"fmt"
	"os"
)

//go:embed shell/cde.zsh
var zshInit string

var shells = map[string]string{
	"zsh": zshInit,
}

func InitShell(shell string) {
	init, ok := shells[shell]
	if !ok {
		fmt.Fprintf(os.Stderr, "unsupported shell: %s\n", shell)
		os.Exit(1)
	}
	fmt.Print(init)
}
