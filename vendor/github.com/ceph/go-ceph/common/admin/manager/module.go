package manager

import (
	"github.com/ceph/go-ceph/internal/commands"
)

// EnableModule will enable the specified manager module.
//
// Similar To:
//
//	ceph mgr module enable <module> [--force]
func (fsa *MgrAdmin) EnableModule(module string, force bool) error {
	m := map[string]string{
		"prefix": "mgr module enable",
		"module": module,
		"format": "json",
	}
	if force {
		m["force"] = "--force"
	}
	// Why is this _only_ part of the mon command json? You'd think a mgr
	// command would be available as a MgrCommand but I couldn't figure it out.
	return commands.MarshalMonCommand(fsa.conn, m).NoData().End()
}

// DisableModule will disable the specified manager module.
//
// Similar To:
//
//	ceph mgr module disable <module>
func (fsa *MgrAdmin) DisableModule(module string) error {
	m := map[string]string{
		"prefix": "mgr module disable",
		"module": module,
		"format": "json",
	}
	return commands.MarshalMonCommand(fsa.conn, m).NoData().End()
}

// DisabledModule describes a disabled Ceph mgr module.
// The Ceph JSON structure contains a complex module_options
// substructure that go-ceph does not currently implement.
type DisabledModule struct {
	Name        string `json:"name"`
	CanRun      bool   `json:"can_run"`
	ErrorString string `json:"error_string"`
}

// ModuleInfo contains fields that report the status of modules within the
// ceph mgr.
type ModuleInfo struct {
	// EnabledModules lists the names of the enabled modules.
	EnabledModules []string `json:"enabled_modules"`
	// AlwaysOnModules lists the names of the always-on modules.
	AlwaysOnModules []string `json:"always_on_modules"`
	// DisabledModules lists structures describing modules that are
	// not currently enabled.
	DisabledModules []DisabledModule `json:"disabled_modules"`
}

func parseModuleInfo(res commands.Response) (*ModuleInfo, error) {
	m := &ModuleInfo{}
	if err := res.NoStatus().Unmarshal(m).End(); err != nil {
		return nil, err
	}
	return m, nil
}

// ListModules returns a module info struct reporting the lists of
// enabled, disabled, and always-on modules in the Ceph mgr.
func (fsa *MgrAdmin) ListModules() (*ModuleInfo, error) {
	m := map[string]string{
		"prefix": "mgr module ls",
		"format": "json",
	}
	return parseModuleInfo(commands.MarshalMonCommand(fsa.conn, m))
}
