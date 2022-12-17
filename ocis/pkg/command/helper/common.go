package helper

import (
	"fmt"
)

// SubcommandDescription
// FIXME: nolint
// nolint: revive
func SubcommandDescription(serviceName string) string {
	return fmt.Sprintf("%s service commands", serviceName)
}
