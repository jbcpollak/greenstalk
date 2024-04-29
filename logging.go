package greenstalk

import (
	"log/slog"

	"github.com/jbcpollak/greenstalk/internal"
)

func SetLogger(logger *slog.Logger) {
	internal.Logger = logger
}
