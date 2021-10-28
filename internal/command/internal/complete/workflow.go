// Package complete provides functions for shell completion of subcommands.
package complete

import (
	"fmt"

	"github.com/tmessi/cci/internal/command/internal/global"
	"github.com/tmessi/cci/internal/command/internal/signal"
	"github.com/tmessi/cci/internal/status"
	"github.com/urfave/cli/v2"
)

// Workflow provides completion values for commands that take a worflow name.
func Workflow(c *cli.Context) {
	ctx, cancel := signal.InitContext()
	defer cancel()

	client, err := global.Client(c)
	if err != nil {
		return
	}

	switch nargs := c.NArg(); {
	case nargs >= 1:
		return
	default:
		s, err := status.Check(ctx, client, c.String("branch"), 1)
		if err != nil {
			return
		}
		if len(s.Pipelines) <= 0 {
			return
		}
		for _, w := range s.Pipelines[0].Workflows {
			fmt.Println(w.Name)
		}
	}
}
