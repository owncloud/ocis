package rados

// #cgo LDFLAGS: -lrados
// #include <stdlib.h>
// #include <rados/librados.h>
import "C"

import (
	"unsafe"

	"github.com/ceph/go-ceph/internal/cutil"
	"github.com/ceph/go-ceph/internal/retry"
)

var argvPlaceholder = "placeholder"

//revive:disable:var-naming old-yet-exported public api

// ClusterStat represents Ceph cluster statistics.
type ClusterStat struct {
	Kb          uint64
	Kb_used     uint64
	Kb_avail    uint64
	Num_objects uint64
}

//revive:enable:var-naming

// Conn is a connection handle to a Ceph cluster.
type Conn struct {
	cluster   C.rados_t
	connected bool
}

// ClusterRef represents a fundamental RADOS cluster connection.
type ClusterRef C.rados_t

// Cluster returns the underlying RADOS cluster reference for this Conn.
func (c *Conn) Cluster() ClusterRef {
	return ClusterRef(c.cluster)
}

// PingMonitor sends a ping to a monitor and returns the reply.
func (c *Conn) PingMonitor(id string) (string, error) {
	cid := C.CString(id)
	defer C.free(unsafe.Pointer(cid))

	var strlen C.size_t
	var strout *C.char

	ret := C.rados_ping_monitor(c.cluster, cid, &strout, &strlen)
	defer C.rados_buffer_free(strout)

	if ret == 0 {
		reply := C.GoStringN(strout, (C.int)(strlen))
		return reply, nil
	}
	return "", getError(ret)
}

// Connect establishes a connection to a RADOS cluster. It returns an error,
// if any.
func (c *Conn) Connect() error {
	ret := C.rados_connect(c.cluster)
	if ret != 0 {
		return getError(ret)
	}
	c.connected = true
	return nil
}

// Shutdown disconnects from the cluster.
func (c *Conn) Shutdown() {
	if err := c.ensureConnected(); err != nil {
		return
	}
	freeConn(c)
}

// ReadConfigFile configures the connection using a Ceph configuration file.
func (c *Conn) ReadConfigFile(path string) error {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	ret := C.rados_conf_read_file(c.cluster, cPath)
	return getError(ret)
}

// ReadDefaultConfigFile configures the connection using a Ceph configuration
// file located at default locations.
func (c *Conn) ReadDefaultConfigFile() error {
	ret := C.rados_conf_read_file(c.cluster, nil)
	return getError(ret)
}

// OpenIOContext creates and returns a new IOContext for the given pool.
//
// Implements:
//  int rados_ioctx_create(rados_t cluster, const char *pool_name,
//                         rados_ioctx_t *ioctx);
func (c *Conn) OpenIOContext(pool string) (*IOContext, error) {
	cPool := C.CString(pool)
	defer C.free(unsafe.Pointer(cPool))
	ioctx := &IOContext{conn: c}
	ret := C.rados_ioctx_create(c.cluster, cPool, &ioctx.ioctx)
	if ret == 0 {
		return ioctx, nil
	}
	return nil, getError(ret)
}

// ListPools returns the names of all existing pools.
func (c *Conn) ListPools() (names []string, err error) {
	buf := make([]byte, 4096)
	for {
		ret := C.rados_pool_list(c.cluster,
			(*C.char)(unsafe.Pointer(&buf[0])), C.size_t(len(buf)))
		if ret < 0 {
			return nil, getError(ret)
		}

		if int(ret) > len(buf) {
			buf = make([]byte, ret)
			continue
		}

		names = cutil.SplitSparseBuffer(buf[:ret])
		return names, nil
	}
}

// SetConfigOption sets the value of the configuration option identified by
// the given name.
func (c *Conn) SetConfigOption(option, value string) error {
	cOpt, cVal := C.CString(option), C.CString(value)
	defer C.free(unsafe.Pointer(cOpt))
	defer C.free(unsafe.Pointer(cVal))
	ret := C.rados_conf_set(c.cluster, cOpt, cVal)
	return getError(ret)
}

// GetConfigOption returns the value of the Ceph configuration option
// identified by the given name.
func (c *Conn) GetConfigOption(name string) (value string, err error) {
	cOption := C.CString(name)
	defer C.free(unsafe.Pointer(cOption))

	var buf []byte
	// range from 4k to 256KiB
	retry.WithSizes(4096, 1<<18, func(size int) retry.Hint {
		buf = make([]byte, size)
		ret := C.rados_conf_get(
			c.cluster,
			cOption,
			(*C.char)(unsafe.Pointer(&buf[0])),
			C.size_t(len(buf)))
		err = getError(ret)
		return retry.DoubleSize.If(err == errNameTooLong)
	})
	if err != nil {
		return "", err
	}
	value = C.GoString((*C.char)(unsafe.Pointer(&buf[0])))
	return value, nil
}

// WaitForLatestOSDMap blocks the caller until the latest OSD map has been
// retrieved.
func (c *Conn) WaitForLatestOSDMap() error {
	ret := C.rados_wait_for_latest_osdmap(c.cluster)
	return getError(ret)
}

func (c *Conn) ensureConnected() error {
	if c.connected {
		return nil
	}
	return ErrNotConnected
}

// GetClusterStats returns statistics about the cluster associated with the
// connection.
func (c *Conn) GetClusterStats() (stat ClusterStat, err error) {
	if err := c.ensureConnected(); err != nil {
		return ClusterStat{}, err
	}
	cStat := C.struct_rados_cluster_stat_t{}
	ret := C.rados_cluster_stat(c.cluster, &cStat)
	if ret < 0 {
		return ClusterStat{}, getError(ret)
	}
	return ClusterStat{
		Kb:          uint64(cStat.kb),
		Kb_used:     uint64(cStat.kb_used),
		Kb_avail:    uint64(cStat.kb_avail),
		Num_objects: uint64(cStat.num_objects),
	}, nil
}

// ParseConfigArgv configures the connection using a unix style command line
// argument vector.
//
// Implements:
//  int rados_conf_parse_argv(rados_t cluster, int argc,
//                            const char **argv);
func (c *Conn) ParseConfigArgv(argv []string) error {
	if c.cluster == nil {
		return ErrNotConnected
	}
	if len(argv) == 0 {
		return ErrEmptyArgument
	}
	cargv := make([]*C.char, len(argv))
	for i := range argv {
		cargv[i] = C.CString(argv[i])
		defer C.free(unsafe.Pointer(cargv[i]))
	}

	ret := C.rados_conf_parse_argv(c.cluster, C.int(len(cargv)), &cargv[0])
	return getError(ret)
}

// ParseCmdLineArgs configures the connection from command line arguments.
//
// This function passes a placeholder value to Ceph as argv[0], see
// ParseConfigArgv for a version of this function that allows the caller to
// specify argv[0].
func (c *Conn) ParseCmdLineArgs(args []string) error {
	argv := make([]string, len(args)+1)
	// Ceph expects a proper argv array as the actual contents with the
	// first element containing the executable name
	argv[0] = argvPlaceholder
	for i := range args {
		argv[i+1] = args[i]
	}
	return c.ParseConfigArgv(argv)
}

// ParseDefaultConfigEnv configures the connection from the default Ceph
// environment variable CEPH_ARGS.
func (c *Conn) ParseDefaultConfigEnv() error {
	ret := C.rados_conf_parse_env(c.cluster, nil)
	return getError(ret)
}

// GetFSID returns the fsid of the cluster as a hexadecimal string. The fsid
// is a unique identifier of an entire Ceph cluster.
func (c *Conn) GetFSID() (fsid string, err error) {
	buf := make([]byte, 37)
	ret := C.rados_cluster_fsid(c.cluster,
		(*C.char)(unsafe.Pointer(&buf[0])), C.size_t(len(buf)))
	// FIXME: the success case isn't documented correctly in librados.h
	if ret == 36 {
		fsid = C.GoString((*C.char)(unsafe.Pointer(&buf[0])))
		return fsid, nil
	}
	return "", getError(ret)
}

// GetInstanceID returns a globally unique identifier for the cluster
// connection instance.
func (c *Conn) GetInstanceID() uint64 {
	// FIXME: are there any error cases for this?
	return uint64(C.rados_get_instance_id(c.cluster))
}

// MakePool creates a new pool with default settings.
func (c *Conn) MakePool(name string) error {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	ret := C.rados_pool_create(c.cluster, cName)
	return getError(ret)
}

// DeletePool deletes a pool and all the data inside the pool.
func (c *Conn) DeletePool(name string) error {
	if err := c.ensureConnected(); err != nil {
		return err
	}
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	ret := C.rados_pool_delete(c.cluster, cName)
	return getError(ret)
}

// GetPoolByName returns the ID of the pool with a given name.
func (c *Conn) GetPoolByName(name string) (int64, error) {
	if err := c.ensureConnected(); err != nil {
		return 0, err
	}
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	ret := int64(C.rados_pool_lookup(c.cluster, cName))
	if ret < 0 {
		return 0, radosError(ret)
	}
	return ret, nil
}

// GetPoolByID returns the name of a pool by a given ID.
func (c *Conn) GetPoolByID(id int64) (string, error) {
	buf := make([]byte, 4096)
	if err := c.ensureConnected(); err != nil {
		return "", err
	}
	cid := C.int64_t(id)
	ret := int(C.rados_pool_reverse_lookup(c.cluster, cid, (*C.char)(unsafe.Pointer(&buf[0])), C.size_t(len(buf))))
	if ret < 0 {
		return "", radosError(ret)
	}
	return C.GoString((*C.char)(unsafe.Pointer(&buf[0]))), nil
}
