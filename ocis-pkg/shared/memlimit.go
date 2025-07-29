package shared

import (
	"log/slog"
	"os"
	"runtime"

	"github.com/KimMachineGun/automemlimit/memlimit"
)

// we init the memlimit here to include it for ocis als well as individual service binaries
func init() {
	slog.SetLogLoggerLevel(slog.LevelError)

	// Enable system memory provider on non-Linux systems to avoid cgroups errors
	if runtime.GOOS != "linux" {
		os.Setenv("AUTOMEMLIMIT_EXPERIMENT", "system")
	}

	_, _ = memlimit.SetGoMemLimitWithOpts(
		memlimit.WithLogger(slog.Default()),
	)
}
