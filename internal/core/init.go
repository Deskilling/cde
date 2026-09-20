package core

import (
	_ "embed"
	"fmt"
	"log/slog"
	"os"
)

//go:embed shell/cde.zsh
var zshInit string

//go:embed shell/cde.fish
var fishInit string

//go:embed shell/cde.bash
var bashInit string

var shells = map[string]string{
	"zsh":  zshInit,
	"fish": fishInit,
	"bash": bashInit,
}

func InitShell(shell string) {
	init, ok := shells[shell]
	if !ok {
		slog.Error("unsupported shell", "shell",shell)
		os.Exit(1)
	}
	fmt.Print(init)
}
