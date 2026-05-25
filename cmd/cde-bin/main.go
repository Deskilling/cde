package main

import (
	"fmt"
	"os"

	"cde/internal/core"
	"cde/internal/editor/zed"
)

func init() {
	core.InitLogger(0)
}

func main() {
	if len(os.Args) == 3 && os.Args[1] == "init" {
		core.InitShell(os.Args[2])
		return
	}

	p, _ := zed.ZedExtractWorkspacePath()
	fmt.Println(p)
}
