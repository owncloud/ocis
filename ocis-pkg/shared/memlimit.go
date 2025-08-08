package shared

import (
	"log/slog"
	"os"

	"github.com/KimMachineGun/automemlimit/memlimit"
)

// we init the memlimit here to include it for ocis als well as individual service binaries
func init() {
	// Check if AUTOMEMLIMIT is set to "off" to disable memory limit functionality
	if os.Getenv("AUTOMEMLIMIT") == "off" {
		return
	}

	slog.SetLogLoggerLevel(slog.LevelError)
	_, _ = memlimit.SetGoMemLimitWithOpts(
		memlimit.WithLogger(slog.Default()),
	)
}
