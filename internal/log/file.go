package log

import (
	"os"
	"path/filepath"

	"github.com/babyname/fate/config"
	"golang.org/x/exp/slog"
)

func New(cfg config.LogConfig) (Logger, error) {
	dir, _ := filepath.Split(cfg.Path)
	if dir != "" {
		_ = os.MkdirAll(dir, 0755)
	}

	file, err := os.OpenFile(cfg.Path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}

	opts.AddSource = cfg.ShowSource
	var h slog.Handler = opts.NewJSONHandler(file)
	if cfg.LogType != "json" {
		h = opts.NewTextHandler(file)
	}

	l := slog.New(h)
	l.Enabled(stringToLevel(cfg.Level))
	return l, nil
}
