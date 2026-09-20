package main

import (
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"

	"cde/internal/core"
	"cde/internal/editor"
)

const Version = "0.0.6"

func init() {
	core.InitLogger(8)

	err := core.LoadConfig(core.GetConfigLocation())
	if err != nil {
		slog.Info("created default config at","location", core.GetConfigLocation())
	}
}

func usage() {
	slog.Info("Usage:")
	fmt.Println("  cde version             shows version")
	fmt.Println("  cde help 		  shows help")
	fmt.Println("  cde install <shell>  	  install automatically for given shell")
	fmt.Println("  cde init <shell>        returns script for given shell")
	fmt.Println("  cde path	          returns latest path")
	fmt.Println("  cde config	          returns config path")
	fmt.Println("  cde list                lists all available editors")
	fmt.Println("  cde -[editor]  	  switch to the latest workspace of specified editor")
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
			if u.Editor.Name() == _editor {
				workspace, err := u.Editor.ExtractWorkspace()
				if err != nil {
					slog.Error("Failed getting workspace", "err", err)
					return
				}
				fmt.Print(workspace.Path)
			}
		}

		return
	}

	switch args[0] {
	case "version":
		fmt.Printf("cde-bin Version %s on %s %s\n", Version, runtime.GOOS, runtime.GOARCH)

	case "path":
		editor.Load()

		w, err := editor.Latest()
		if err != nil {
			slog.Error("Failed returning path", "err", err)
			return
		}
		fmt.Print(w.Path)

	case "init":
		if len(args) != 2 {
			slog.Error("init requires a single <shell> argument")
			usage()
			return
		}

		core.InitShell(args[1])

	case "install":
		if len(args) != 2 {
			slog.Error("install requires a single <shell> argument")
			usage()
			return
		}

		core.InstallShell(args[1])

	case "config":
		fmt.Print(core.GetConfigLocation())

	case "list":
		editor.Load()

		for _, v := range editor.Registered {
			fmt.Println(v.Editor.Name())
		}

	default:
		slog.Error("unknown command", "arg",args[0])
		slog.Info("see all valid arguments via cde help")
		return
	}

}
