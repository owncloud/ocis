package admin

// For APIs that accept extra sets of "boolean" flags we may end up wanting
// multiple different sets of supported flags. Example: most rm functions
// accept a force flag, but only subvolume delete has retain snapshots.
// To make this somewhat uniform in the admin package we define a utility
// interface and helper function to merge flags with naming options.

type flagSet interface {
	flags() map[string]bool
}

type commonRmFlags struct {
	force bool
}

func (f commonRmFlags) flags() map[string]bool {
	o := make(map[string]bool)
	if f.force {
		o["force"] = true
	}
	return o
}

// SubVolRmFlags does not embed other types to simplify and keep the
// interface with the type flat and simple. At the cost of some code
// duplication we get a nicer UX for those using the library.

// SubVolRmFlags may be used to specify behavior modifying flags when
// removing sub volumes.
type SubVolRmFlags struct {
	Force           bool
	RetainSnapshots bool
}

func (f SubVolRmFlags) flags() map[string]bool {
	o := make(map[string]bool)
	if f.Force {
		o["force"] = true
	}
	if f.RetainSnapshots {
		o["retain_snapshots"] = true
	}
	return o
}

// mergeFlags combines a set of key-value settings with any type implementing
// the flagSet interface.
func mergeFlags(m map[string]string, f flagSet) map[string]interface{} {
	o := make(map[string]interface{})
	for k, v := range m {
		o[k] = v
	}
	for k, v := range f.flags() {
		o[k] = v
	}
	return o
}
