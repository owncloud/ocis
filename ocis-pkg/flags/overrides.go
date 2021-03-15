package flags

// OverrideDefaultString checks whether the default value of v is the zero value, if so, ensure the flag has a correct
// value by providing one. A value different than zero would mean that it was read from a config file either from an
// extension or from a higher source (i.e: ocis command).
func OverrideDefaultString(v, def string) string {
	if v != "" {
		return v
	}

	return def
}

// OverrideDefaultBool checks whether the default value of v is the zero value, if so, ensure the flag has a correct
// value by providing one. A value different than zero would mean that it was read from a config file either from an
// extension or from a higher source (i.e: ocis command).
func OverrideDefaultBool(v, def bool) bool {
	if v {
		return v
	}

	return def
}

// OverrideDefaultInt checks whether the default value of v is the zero value, if so, ensure the flag has a correct
// value by providing one. A value different than zero would mean that it was read from a config file either from an
// extension or from a higher source (i.e: ocis command).
func OverrideDefaultInt(v, def int) int {
	if v != 0 {
		return v
	}

	return def
}

// OverrideDefaultInt64 checks whether the default value of v is the zero value, if so, ensure the flag has a correct
// value by providing one. A value different than zero would mean that it was read from a config file either from an
// extension or from a higher source (i.e: ocis command).
func OverrideDefaultInt64(v, def int64) int64 {
	if v != 0 {
		return v
	}

	return def
}

// OverrideDefaultUint64 checks whether the default value of v is the zero value, if so, ensure the flag has a correct
// value by providing one. A value different than zero would mean that it was read from a config file either from an
// extension or from a higher source (i.e: ocis command).
func OverrideDefaultUint64(v, def uint64) uint64 {
	if v != 0 {
		return v
	}

	return def
}
