package fate

import (
	"log/slog"
)

var log *slog.Logger

func init() {
	log = slog.Default()
}
