package defaults

import (
	"log"
	"os"
	"path"
)

var (
	// BaseDataPathType
	// switch between modes
	// FIXME: nolint
	// nolint: revive
	BaseDataPathType = "homedir" // or "path"
	// BaseDataPathValue
	// default data path
	// FIXME: nolint
	// nolint: revive
	BaseDataPathValue = "/var/lib/ocis"
)

// BaseDataPath
// FIXME: nolint
// nolint: revive
func BaseDataPath() string {

	// It is not nice to have hidden / secrete configuration options
	// But how can we update the base path for every occurrence with a flagset option?
	// This is currently not possible and needs a new configuration concept
	p := os.Getenv("OCIS_BASE_DATA_PATH")
	if p != "" {
		return p
	}

	switch BaseDataPathType {
	case "homedir":
		dir, err := os.UserHomeDir()
		if err != nil {
			// fallback to BaseDatapathValue for users without home
			return BaseDataPathValue
		}
		return path.Join(dir, ".ocis")
	case "path":
		return BaseDataPathValue
	default:
		log.Fatalf("BaseDataPathType %s not found", BaseDataPathType)
		return ""
	}
}

var (
	// BaseConfigPathType
	// switch between modes
	// FIXME: nolint
	// nolint: revive
	BaseConfigPathType = "homedir" // or "path"
	// BaseConfigPathValue
	// default config path
	// FIXME: nolint
	// nolint: revive
	BaseConfigPathValue = "/etc/ocis"
)

// BaseConfigPath
// FIXME: nolint
// nolint: revive
func BaseConfigPath() string {
	// It is not nice to have hidden / secrete configuration options
	// But how can we update the base path for every occurrence with a flagset option?
	// This is currently not possible and needs a new configuration concept
	p := os.Getenv("OCIS_CONFIG_DIR")
	if p != "" {
		return p
	}

	switch BaseConfigPathType {
	case "homedir":
		dir, err := os.UserHomeDir()
		if err != nil {
			// fallback to BaseConfigPathValue for users without home
			return BaseConfigPathValue
		}
		return path.Join(dir, ".ocis", "config")
	case "path":
		return BaseConfigPathValue
	default:
		log.Fatalf("BaseConfigPathType %s not found", BaseConfigPathType)
		return ""
	}
}
