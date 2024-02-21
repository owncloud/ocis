package useragent

import (
	"fmt"
	"strconv"
	"strings"
)

type VersionNo struct {
	Major int
	Minor int
	Patch int
}

func parseVersion(ver string, verno *VersionNo) {
	var err error
	parts := strings.Split(ver, ".")
	if len(parts) > 0 {
		if verno.Major, err = strconv.Atoi(parts[0]); err != nil {
			return
		}
	}
	if len(parts) > 1 {
		if verno.Minor, err = strconv.Atoi(parts[1]); err != nil {
			return
		}
		if len(parts) > 2 {
			if verno.Patch, err = strconv.Atoi(parts[2]); err != nil {
				return
			}
		}
	}
}

// VersionNoShort return version string in format <Major>.<Minor>
func (ua UserAgent) VersionNoShort() string {
	if ua.VersionNo.Major == 0 && ua.VersionNo.Minor == 0 && ua.VersionNo.Patch == 0 {
		return ""
	}
	return fmt.Sprintf("%d.%d", ua.VersionNo.Major, ua.VersionNo.Minor)
}

// VersionNoFull returns version string in format <Major>.<Minor>.<Patch>
func (ua UserAgent) VersionNoFull() string {
	if ua.VersionNo.Major == 0 && ua.VersionNo.Minor == 0 && ua.VersionNo.Patch == 0 {
		return ""
	}
	return fmt.Sprintf("%d.%d.%d", ua.VersionNo.Major, ua.VersionNo.Minor, ua.VersionNo.Patch)
}

// OSVersionNoShort returns OS version string in format <Major>.<Minor>
func (ua UserAgent) OSVersionNoShort() string {
	if ua.OSVersionNo.Major == 0 && ua.OSVersionNo.Minor == 0 && ua.OSVersionNo.Patch == 0 {
		return ""
	}
	return fmt.Sprintf("%d.%d", ua.OSVersionNo.Major, ua.OSVersionNo.Minor)
}

// OSVersionNoFull returns OS version string in format <Major>.<Minor>.<Patch>
func (ua UserAgent) OSVersionNoFull() string {
	if ua.OSVersionNo.Major == 0 && ua.OSVersionNo.Minor == 0 && ua.OSVersionNo.Patch == 0 {
		return ""
	}
	return fmt.Sprintf("%d.%d.%d", ua.OSVersionNo.Major, ua.OSVersionNo.Minor, ua.OSVersionNo.Patch)
}
