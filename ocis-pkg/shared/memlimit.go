package shared

import (
	"log/slog"
	"os"

	"github.com/KimMachineGun/automemlimit/memlimit"
)

// we init the memlimit here to include it for ocis als well as individual service binaries
func init() {
	slog.SetLogLoggerLevel(slog.LevelError)

	if os.Getenv("AUTOMEMLIMIT") == "off" {
		return
	}

	_, _ = memlimit.SetGoMemLimitWithOpts(
		memlimit.WithLogger(slog.Default()),
	)
}
