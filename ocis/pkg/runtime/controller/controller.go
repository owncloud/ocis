package controller

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/owncloud/ocis/ocis/pkg/runtime/config"
	"github.com/owncloud/ocis/ocis/pkg/runtime/process"
	"github.com/owncloud/ocis/ocis/pkg/runtime/storage"
	"github.com/owncloud/ocis/ocis/pkg/runtime/watcher"
	"github.com/rs/zerolog"

	"github.com/olekukonko/tablewriter"
)

// Controller supervises processes.
type Controller struct {
	m       *sync.RWMutex
	options Options
	log     zerolog.Logger
	Config  *config.Config

	Store storage.Storage

	// Bin is the oCIS single binary name.
	Bin string

	// BinPath is the oCIS single binary path withing the host machine.
	// The Controller needs to know the binary location in order to spawn new extensions.
	BinPath string

	// Terminated facilitates communication from Watcher <-> Controller. Writes to this
	// channel WILL always attempt to restart the crashed process.
	Terminated chan process.ProcEntry
}

var (
	once = sync.Once{}
)

// NewController initializes a new controller.
func NewController(o ...Option) Controller {
	opts := &Options{}

	for _, f := range o {
		f(opts)
	}

	c := Controller{
		m:          &sync.RWMutex{},
		options:    *opts,
		log:        *opts.Log,
		Bin:        "ocis",
		Terminated: make(chan process.ProcEntry),
		Store:      storage.NewMapStorage(),

		Config: opts.Config,
	}

	if opts.Bin != "" {
		c.Bin = opts.Bin
	}

	// Get binary location from $PATH lookup. If not present, it uses arg[0] as entry point.
	path, err := exec.LookPath(c.Bin)
	if err != nil {
		c.log.Debug().Msg("oCIS binary not present in PATH, using Args[0]")
		path = os.Args[0]
	}
	c.BinPath = path
	return c
}

// Start and watches a process.
func (c *Controller) Start(pe process.ProcEntry) error {
	if pid := c.Store.Load(pe.Extension); pid != 0 {
		c.log.Debug().Msg(fmt.Sprintf("extension already running: %s", pe.Extension))
		return nil
	}

	w := watcher.NewWatcher()
	if err := pe.Start(c.BinPath); err != nil {
		return err
	}

	// store the spawned child process PID.
	if err := c.Store.Store(pe); err != nil {
		return err
	}

	w.Follow(pe, c.Terminated, c.options.Config.KeepAlive)

	once.Do(func() {
		j := janitor{
			time.Second,
			c.Store,
		}

		go j.run()
		go detach(c)
	})
	return nil
}

// Kill a managed process.
// Should a process managed by the runtime be allowed to be killed if the runtime is configured not to?
func (c *Controller) Kill(pe process.ProcEntry) error {
	// load stored PID
	pid := c.Store.Load(pe.Extension)

	// find process in host by PID
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	if err := c.Store.Delete(pe); err != nil {
		return err
	}
	c.log.Info().Str("package", "watcher").Msgf("terminating %v", pe.Extension)

	// terminate child process
	return p.Kill()
}

// Shutdown a running runtime.
func (c *Controller) Shutdown(ch chan struct{}) error {
	entries := c.Store.LoadAll()
	for cmd, pid := range entries {
		c.log.Info().Str("package", "watcher").Msgf("gracefully terminating %v", cmd)
		p, _ := os.FindProcess(pid)
		if err := p.Kill(); err != nil {
			return err
		}
	}

	ch <- struct{}{}
	return nil
}

// List managed processes.
func (c *Controller) List() string {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetHeader([]string{"Extension", "PID"})

	entries := c.Store.LoadAll()

	keys := make([]string, 0, len(entries))
	for k := range entries {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, v := range keys {
		table.Append([]string{v, strconv.Itoa(entries[v])})
	}

	table.Render()
	return tableString.String()
}
