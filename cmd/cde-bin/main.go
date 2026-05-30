package main

import (
	"fmt"
	"os"

	"cde/internal/core"
	"cde/internal/editor"

	"charm.land/log/v2"
)

func init() {
	core.InitLogger(0)
	core.LoadConfig(core.GetConfigLocation())
	editor.Load()
}

func usage() {
	log.Print("Usage:")
	log.Print("  cde help 			   shows help")
	log.Print("  cde install <shell>	 install automatically for your shell")
	log.Print("  cde init <shell>		returns script for current shell")
	log.Print("  cde path				return latest path")
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 || args[0] == "help" {
		usage()
		return
	}

	switch args[0] {
	case "path":
		w, err := editor.Latest()
		if err != nil {
			log.Error(err)
			return
		}
		fmt.Print(w.Path)

	case "init":
		if len(args) != 2 {
			log.Error("init requires a single <shell> argument")
			usage()
			return
		}

		core.InitShell(args[1])

	case "install":
		if len(args) != 2 {
			log.Error("install requires a single <shell> argument")
			usage()
			return
		}

		core.InstallShell(args[1])

	default:
		log.Errorf("unknown command: %s", args[0])
		log.Info("see all valid arguments via cde help")
		return
	}

}
