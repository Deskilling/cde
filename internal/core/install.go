package core

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func appendToFile(path, content string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open %s: %w", path, err)
	}
	defer f.Close()

	_, err = fmt.Fprintf(f, "\n%s\n", content)
	if err != nil {
		return fmt.Errorf("failed to write to %s: %w", path, err)
	}

	return nil
}

func InstallShell(shell string) error {
	switch shell {
	case "zsh":
		return appendToFile(filepath.Join(os.Getenv("HOME"), ".zshrc"), `eval "$(cde-bin init zsh)"`)
	default:
		return errors.New("shell not supported")
	}
}
