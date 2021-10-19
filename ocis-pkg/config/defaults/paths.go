package defaults

import (
	"log"
	"os"
	"path"
)

const ()

var (
	// switch between modes
	BaseDataPathType = "homedir"
	// don't read from this, only write
	BaseDataPathValue = "/var/lib/ocis"
)

func BaseDataPath() string {

	// It is not nice to have hidden / secrete configuration options
	// But how can we update the base path for every occurence with a flageset option?
	// This is currenlty not possible and needs a new configuration concept
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
