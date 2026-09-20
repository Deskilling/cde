package core

import (
	"fmt"
	"os"
	"runtime"
)

func supported(paths map[string]string) (string, error) {
	path, ok := paths[runtime.GOOS]
	if !ok {
		return "", fmt.Errorf("unsupported os: %s", runtime.GOOS)
	}

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("file doest not exists: %w", err)
	}
	return path, nil
}

// cheks if config overrites default path
// also checks if the files exists
func CheckOverride(editorName string, defaultPaths map[string]string) (string, error) {
	editor, ok := GetConfig().Editors[editorName]
	if editor.WorkspacePath[runtime.GOOS] == "" || !ok {
		return supported(defaultPaths)
	}
	return supported(editor.WorkspacePath)
}
