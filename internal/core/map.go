package core

import (
	"fmt"
	"runtime"
)

func supported(paths map[string]string) (string, error) {
	path, ok := paths[runtime.GOOS]
	if !ok {
		return "", fmt.Errorf("unsupported os: %s", runtime.GOOS)
	}
	return path, nil
}

func CheckOverride(editorName string, defaultPaths map[string]string) (string, error) {
	editor, ok := GetConfig().Editors[editorName]
	if editor.WorkspacePath[runtime.GOOS] == "" || !ok {
		return supported(defaultPaths)
	}
	return supported(editor.WorkspacePath)
}
