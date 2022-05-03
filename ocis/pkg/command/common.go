package command

import (
	"fmt"
)

func subcommandDescription(serviceName string) string {
	return fmt.Sprintf("%s extension commands", serviceName)
}
