package admin

import (
	"strconv"

	ccom "github.com/ceph/go-ceph/common/commands"
	"github.com/ceph/go-ceph/internal/commands"
	"github.com/ceph/go-ceph/rados"
)

// RadosCommander provides an interface to execute JSON-formatted commands that
// allow the cephfs administrative functions to interact with the Ceph cluster.
type RadosCommander = ccom.RadosCommander

// FSAdmin is used to administrate CephFS within a ceph cluster.
type FSAdmin struct {
	conn RadosCommander
}

// New creates an FSAdmin automatically based on the default ceph
// configuration file. If more customization is needed, create a
// *rados.Conn as you see fit and use NewFromConn to use that
// connection with these administrative functions.
func New() (*FSAdmin, error) {
	conn, err := rados.NewConn()
	if err != nil {
		return nil, err
	}
	err = conn.ReadDefaultConfigFile()
	if err != nil {
		return nil, err
	}
	err = conn.Connect()
	if err != nil {
		return nil, err
	}
	return NewFromConn(conn), nil
}

// NewFromConn creates an FSAdmin management object from a preexisting
// rados connection. The existing connection can be rados.Conn or any
// type implementing the RadosCommander interface. This may be useful
// if the calling layer needs to inject additional logging, error handling,
// fault injection, etc.
func NewFromConn(conn RadosCommander) *FSAdmin {
	return &FSAdmin{conn}
}

func (fsa *FSAdmin) validate() error {
	if fsa.conn == nil {
		return rados.ErrNotConnected
	}
	return nil
}

// rawMgrCommand takes a byte buffer and sends it to the MGR as a command.
// The buffer is expected to contain preformatted JSON.
func (fsa *FSAdmin) rawMgrCommand(buf []byte) response {
	return commands.RawMgrCommand(fsa.conn, buf)
}

// marshalMgrCommand takes an generic interface{} value, converts it to JSON and
// sends the json to the MGR as a command.
func (fsa *FSAdmin) marshalMgrCommand(v interface{}) response {
	return commands.MarshalMgrCommand(fsa.conn, v)
}

// rawMonCommand takes a byte buffer and sends it to the MON as a command.
// The buffer is expected to contain preformatted JSON.
func (fsa *FSAdmin) rawMonCommand(buf []byte) response {
	return commands.RawMonCommand(fsa.conn, buf)
}

// marshalMonCommand takes an generic interface{} value, converts it to JSON and
// sends the json to the MGR as a command.
func (fsa *FSAdmin) marshalMonCommand(v interface{}) response {
	return commands.MarshalMonCommand(fsa.conn, v)
}

type listNamedResult struct {
	Name string `json:"name"`
}

func parseListNames(res response) ([]string, error) {
	var r []listNamedResult
	if err := res.NoStatus().Unmarshal(&r).End(); err != nil {
		return nil, err
	}
	vl := make([]string, len(r))
	for i := range r {
		vl[i] = r[i].Name
	}
	return vl, nil
}

func parseListKeyValues(res response) (map[string]string, error) {
	var x map[string]string
	if err := res.NoStatus().Unmarshal(&x).End(); err != nil {
		return nil, err
	}

	return x, nil
}

// parsePathResponse returns a cleaned up path from requests that get a path
// unless an error is encountered, then an error is returned.
func parsePathResponse(res response) (string, error) {
	if res2 := res.NoStatus(); !res2.Ok() {
		return "", res.End()
	}
	b := res.Body()
	// if there's a trailing newline in the buffer strip it.
	// ceph assumes a CLI wants the output of the buffer and there's
	// no format=json mode available currently.
	for len(b) >= 1 && b[len(b)-1] == '\n' {
		b = b[:len(b)-1]
	}
	return string(b), nil
}

// modeString converts a unix-style mode value to a string-ified version in an
// octal representation (e.g. "777", "700", etc). This format is expected by
// some of the ceph JSON command inputs.
func modeString(m int, force bool) string {
	if force || m != 0 {
		return strconv.FormatInt(int64(m), 8)
	}
	return ""
}

// uint64String converts a uint64 to a string. Some of the ceph json commands
// can take a string or "int" (as a string). This is a common function for
// doing that conversion.
func uint64String(v uint64) string {
	return strconv.FormatUint(uint64(v), 10)
}
