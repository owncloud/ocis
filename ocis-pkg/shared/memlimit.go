package shared

import (
	"log/slog"
	"runtime"

	"github.com/KimMachineGun/automemlimit/memlimit"
)

// we init the memlimit here to include it for ocis als well as individual service binaries
func init() {
	slog.SetLogLoggerLevel(slog.LevelError)

	opts := []memlimit.Option{
		memlimit.WithLogger(slog.Default()),
	}

	// Use the system memory provider on non-Linux systems to avoid cgroups errors
	if runtime.GOOS != "linux" {
		opts = append(opts, memlimit.WithProvider(memlimit.FromSystem))
	}

	_, _ = memlimit.Set(opts...)
}
