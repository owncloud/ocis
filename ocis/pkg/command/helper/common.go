package helper

import (
	"fmt"
)

func SubcommandDescription(serviceName string) string {
	return fmt.Sprintf("%s service commands", serviceName)
}
