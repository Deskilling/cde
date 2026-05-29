package core

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
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

var zshConfig map[string]string = map[string]string{
	"darwin": filepath.Join(os.Getenv("HOME"), ".zshrc"),
	"linux":  filepath.Join(os.Getenv("HOME"), ".zshrc"),
}

var fishConfig map[string]string = map[string]string{
	"darwin": filepath.Join(os.Getenv("HOME"), ".config", "fish", "conf.d"),
	"linux":  filepath.Join(os.Getenv("HOME"), ".config", "fish", "conf.d"),
}

func InstallShell(shell string) (err error) {
	switch shell {
	case "zsh":
		zshConfig, ok := zshConfig[runtime.GOOS]
		if !ok {
			return fmt.Errorf("unsupported os: %s", runtime.GOOS)
		}
		return appendToFile(zshConfig, `eval "$(cde-bin init zsh)"`)

	case "fish":
		confDir, ok := fishConfig[runtime.GOOS]
		if !ok {
			return fmt.Errorf("unsupported os: %s", runtime.GOOS)
		}
		err = os.MkdirAll(confDir, 0755)
		if err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(confDir, "cde.fish"), []byte(`cde-bin init fish | source`), 0644)

	default:
		return errors.New("shell not supported")
	}
}
