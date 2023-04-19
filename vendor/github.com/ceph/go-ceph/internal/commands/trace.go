package commands

import (
	"fmt"

	ccom "github.com/ceph/go-ceph/common/commands"
)

// NewTraceCommander is a RadosCommander that wraps a given RadosCommander
// and when commands are executes prints debug level "traces" to the
// standard output.
func NewTraceCommander(c ccom.RadosCommander) ccom.RadosCommander {
	return &tracingCommander{c}
}

// tracingCommander serves two purposes: first, it allows one to trace the
// input and output json when running the tests. It can help with actually
// debugging the tests. Second, it demonstrates the rationale for using an
// interface in FSAdmin. You can layer any sort of debugging, error injection,
// or whatnot between the FSAdmin layer and the RADOS layer.
type tracingCommander struct {
	conn ccom.RadosCommander
}

func (t *tracingCommander) MgrCommand(buf [][]byte) ([]byte, string, error) {
	fmt.Println("(MGR Command)")
	for i := range buf {
		fmt.Println("IN:", string(buf[i]))
	}
	r, s, err := t.conn.MgrCommand(buf)
	fmt.Println("OUT(result):", string(r))
	if s != "" {
		fmt.Println("OUT(status):", s)
	}
	if err != nil {
		fmt.Println("OUT(error):", err.Error())
	}
	return r, s, err
}

func (t *tracingCommander) MonCommand(buf []byte) ([]byte, string, error) {
	fmt.Println("(MON Command)")
	fmt.Println("IN:", string(buf))
	r, s, err := t.conn.MonCommand(buf)
	fmt.Println("OUT(result):", string(r))
	if s != "" {
		fmt.Println("OUT(status):", s)
	}
	if err != nil {
		fmt.Println("OUT(error):", err.Error())
	}
	return r, s, err
}
