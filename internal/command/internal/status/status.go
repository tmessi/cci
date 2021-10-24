// Package status provides the status subcommand.
package status

import (
	"fmt"

	"github.com/tmessi/cci/internal/command/internal/global"
	"github.com/tmessi/cci/internal/command/internal/signal"
	"github.com/tmessi/cci/internal/status"
	"gopkg.in/urfave/cli.v1"
)

// Command is the status subcommand.
var Command = cli.Command{
	Name:    "status",
	Aliases: []string{"branch-status", "s"},
	Usage:   "Show the status of a branch",
	Action:  action,
}

func action(c *cli.Context) error {
	ctx, cancel := signal.InitContext()
	defer cancel()

	client, err := global.Client(c)
	if err != nil {
		return cli.NewExitError(err.Error(), -1)
	}

	s, err := status.Check(ctx, client, c.GlobalString("branch"))
	if err != nil {
		return cli.NewExitError(err.Error(), -1)
	}

	fmt.Println(s)
	return nil
}
