package main

import (
	"fmt"
	"os"

	"cde/internal/core"
	"cde/internal/editor"

	"charm.land/log/v2"
)

func init() {
	core.InitLogger(-4)
	editor.Load()
}

func main() {
	if len(os.Args) == 3 && os.Args[1] == "init" {
		core.InitShell(os.Args[2])
		return
	}

	w, err := editor.Latest()
	if err != nil {
		log.Error(err)
		return
	}

	fmt.Print(w.Path)

}
