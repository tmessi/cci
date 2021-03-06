// Package complete provides functions for shell completion of subcommands.
package complete

import (
	"fmt"
	"strconv"

	"github.com/tmessi/cci/internal/command/internal/global"
	"github.com/tmessi/cci/internal/command/internal/signal"
	"github.com/tmessi/cci/internal/status"
	"github.com/urfave/cli/v2"
)

// Job provides completion values for commands that take Job numbers or
// worflow and job names.
func Job(c *cli.Context) {
	ctx, cancel := signal.InitContext()
	defer cancel()

	client, err := global.Client(c)
	if err != nil {
		return
	}

	switch nargs := c.NArg(); {
	case nargs >= 2:
		return
	case nargs == 1:
		// If first are is a Job number, no additional args are needed.
		if _, err := strconv.Atoi(c.Args().Get(0)); err == nil {
			return
		}

		workflowName := c.Args().Get(0)
		s, err := status.Check(ctx, client, c.String("branch"), 1)
		if err != nil {
			return
		}
		workflow := s.Workflow(workflowName)
		if workflow == nil {
			return
		}

		for _, j := range workflow.Jobs {
			fmt.Println(j.Name)
		}
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
			for _, j := range w.Jobs {
				fmt.Println(j.Number)
			}
		}
	}
}
