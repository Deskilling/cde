package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"cde/internal/core"
	"cde/internal/editor"

	"charm.land/log/v2"
)

const Version = "0.0.4"

func init() {
	core.InitLogger(log.InfoLevel)

	err := core.LoadConfig(core.GetConfigLocation())
	if err != nil {
		log.Infof("created default config at %s", core.GetConfigLocation())
	}
	editor.Load()
}

func usage() {
	log.Print("Usage:")
	log.Print("  cde version             shows version")
	log.Print("  cde help 			   shows help")
	log.Print("  cde install <shell>	 install automatically for given shell")
	log.Print("  cde init <shell>		returns script for given shell")
	log.Print("  cde path				returns latest path")
	log.Print("  cde config			  returns config path")
	log.Print("  cde list 			   lists all available editors")
	log.Print("  cde -[editor]  		 switch to the latest workspace of specified editor")
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 || args[0] == "help" {
		usage()
		return
	}

	_editor, hasDash := strings.CutPrefix(args[0], "-")
	if hasDash {
		for _, u := range editor.Registered {
			if u.Name() == _editor {
				workspace, err := u.ExtractWorkspace()
				if err != nil {
					log.Error(err)
					return
				}
				fmt.Print(workspace.Path)
			}
		}

		return
	}

	switch args[0] {
	case "version":
		log.Printf("cde-bin Version %s on %s %s", Version, runtime.GOOS, runtime.GOARCH)

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

	case "config":
		fmt.Print(core.GetConfigLocation())

	case "list":
		for _, v := range editor.Registered {
			fmt.Println(v.Name())
		}

	default:
		log.Errorf("unknown command: %s", args[0])
		log.Info("see all valid arguments via cde help")
		return
	}

}
