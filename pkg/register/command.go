package register

import (
	"context"

	"github.com/micro/cli"
	"github.com/oklog/run"
	"github.com/owncloud/ocis/pkg/config"
)

var (
	// Commands defines the slice of commands.
	Commands = []Command{}

	// Handlers defines the slice of handlers.
	Handlers = []Handler{}
)

// Command defines the register command.
type Command func(*config.Config) cli.Command

// Handler defines the register handler.
type Handler func(context.Context, context.CancelFunc, *run.Group, *config.Config) error

// AddCommand appends a command to Commands.
func AddCommand(cmd Command) {
	Commands = append(
		Commands,
		cmd,
	)
}

// AddHandler appends a handler to Handlers.
func AddHandler(hdl Handler) {
	Handlers = append(
		Handlers,
		hdl,
	)
}
