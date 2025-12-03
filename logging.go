package greenstalk

import (
	"log/slog"

	"github.com/jbcpollak/greenstalk/v2/internal"
)

func SetLogger(logger *slog.Logger) {
	internal.Logger = logger
}
