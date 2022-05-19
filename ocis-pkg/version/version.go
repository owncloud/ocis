package version

import (
	"time"

	"github.com/Masterminds/semver"
)

var (
	// String gets defined by the build system.
	String = "dev"

	// Date indicates the build date.
	Date = time.Now().Format("20060102")

	// Legacy defines the old long 4 number ownCloud version needed for some clients
	Legacy = "10.11.0.0"

	// LegacyString defines the old ownCloud version needed for some clients
	LegacyString = "10.11.0"
)

// Compiled returns the compile time of this service.
func Compiled() time.Time {
	t, _ := time.Parse("20060102", Date)
	return t
}

// GetString returns a version string with pre-releases and metadata
func GetString() string {
	return Parsed().String()
}

// Parsed returns a semver Version
func Parsed() *semver.Version {
	versionToParse := String
	if String == "dev" {
		versionToParse = "0.0.0+dev"
	}
	parsedVersion, err := semver.NewVersion(versionToParse)
	// We have no semver version but a commitid
	if err != nil {
		parsedVersion, err = semver.NewVersion("0.0.0+" + String)
		// this should never happen
		if err != nil {
			return &semver.Version{}
		}
	}
	return parsedVersion
}

// ParsedLegacy returns the legacy version
func ParsedLegacy() *semver.Version {
	parsedVersion, err := semver.NewVersion(LegacyString)
	if err != nil {
		return &semver.Version{}
	}
	return parsedVersion
}
