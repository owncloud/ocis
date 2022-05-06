package version

import (
	"strconv"
	"time"

	"github.com/Masterminds/semver"
)

var (
	// String gets defined by the build system.
	String = "dev"

	// Date indicates the build date.
	Date = time.Now().Format("20060102")
)

// Compiled returns the compile time of this service.
func Compiled() time.Time {
	t, _ := time.Parse("20060102", Date)
	return t
}

func GetString() string {
	if String == "dev" {
		return "0.0.0+dev"
	}
	parsedVersion, err := semver.NewVersion(String)
	// We have no semver version but a commitid
	if err != nil {
		return String
	}
	return parsedVersion.String()
}

func Parsed() *semver.Version {
	var versionToParse string
	if String == "dev" {
		versionToParse = "0.0.0+dev"
	}
	parsedVersion, err := semver.NewVersion(versionToParse)
	// We have no semver version but a commitid
	if err != nil {
		parsedVersion, _ = semver.NewVersion("0.0.0+" + String)
	}
	return parsedVersion
}

func Long() string {
	var s string
	if Parsed().Metadata() == "" {
		s = "-" + Parsed().Prerelease()
	}
	s = "+" + Parsed().Metadata()
	return strconv.FormatInt(Parsed().Major(), 10) + "." +
		strconv.FormatInt(Parsed().Minor(), 10) + "." +
		strconv.FormatInt(Parsed().Patch(), 10) + "." + "0" + s
}
