package core

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

func InitLogger(level int) {
	isDebug := false
	if level == -4 {
		isDebug = true
	}

	slog.SetDefault(slog.New(tint.NewTextHandler(os.Stderr, &tint.Options{
		AddSource: isDebug,
		Level:     slog.Level(level),
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey && len(groups) == 0 {
				return slog.Attr{}
			}
			return a
		},
	})))
}
