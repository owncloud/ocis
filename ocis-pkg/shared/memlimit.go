package shared

import (
	"log/slog"

	"github.com/KimMachineGun/automemlimit/memlimit"
)

// we init the memlimit here to include it for ocis als well as individual service binaries
func init() {
	slog.SetLogLoggerLevel(slog.LevelError)
	_, _ = memlimit.SetGoMemLimitWithOpts(
		memlimit.WithLogger(slog.Default()),
	)
}
