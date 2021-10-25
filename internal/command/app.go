// Package command defines the cli.
package command

import (
	"github.com/urfave/cli/v2"

	"github.com/tmessi/cci/internal/command/internal/global"
	"github.com/tmessi/cci/internal/command/internal/output"
	"github.com/tmessi/cci/internal/command/internal/status"
)

// App returns the cli.App with its subcommands and flags.
func App() *cli.App {
	app := cli.NewApp()
	app.Name = "cci"
	app.Usage = "get status and output from circleci"
	app.EnableBashCompletion = true
	app.Version = "0.1"

	app.Flags = global.Flags
	app.Action = status.Command.Action
	app.Commands = []*cli.Command{
		status.Command,
		output.Command,
	}

	return app
}
