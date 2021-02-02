// +build !windows

package process

import (
	"os"

	sys "golang.org/x/sys/unix"
)

// ProcEntry is an entry in the File db.
type ProcEntry struct {
	Args      []string
	Env       []string
	Pid       int
	Extension string
}

// NewProcEntry returns a new ProcEntry.
func NewProcEntry(extension string, env []string, args ...string) ProcEntry {
	return ProcEntry{
		Extension: extension,
		Args:      args,
		Env:       env,
	}
}

// Start a process.
func (e *ProcEntry) Start(binPath string) error {
	var argv = []string{binPath}
	argv = append(argv, e.Args...)

	p, err := os.StartProcess(binPath, argv, &os.ProcAttr{
		Files: []*os.File{
			os.Stdin,
			os.Stdout,
			os.Stderr,
		},
		Env: e.Env,
		Sys: &sys.SysProcAttr{
			Setpgid: true,
		},
	})
	if err != nil {
		return err
	}
	e.Pid = p.Pid

	return nil
}

// Kill the wrapped process.
func (e *ProcEntry) Kill() error {
	p, err := os.FindProcess(e.Pid)
	if err != nil {
		return err
	}

	return p.Kill()
}
