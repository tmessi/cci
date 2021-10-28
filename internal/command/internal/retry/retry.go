// Package retry provides the retry subcommand.
package retry

import (
	"fmt"

	"github.com/tmessi/cci/internal/command/internal/complete"
	"github.com/tmessi/cci/internal/command/internal/global"
	"github.com/tmessi/cci/internal/command/internal/signal"
	"github.com/tmessi/cci/internal/retry"
	"github.com/tmessi/cci/internal/status"
	"github.com/urfave/cli/v2"
)

// Command is the retry subcommand.
var Command = &cli.Command{
	Name:         "retry",
	ArgsUsage:    "<workflow name>",
	Aliases:      []string{"r"},
	Usage:        "Retry a build",
	BashComplete: complete.Workflow,
	Action:       action,
}

func action(c *cli.Context) error {
	ctx, cancel := signal.InitContext()
	defer cancel()

	client, err := global.Client(c)
	if err != nil {
		return cli.NewExitError(err.Error(), -1)
	}

	var workflowID string

	switch c.NArg() {
	case 1:
		workflowName := c.Args().Get(0)

		s, err := status.Check(ctx, client, c.String("branch"))
		if err != nil {
			return cli.NewExitError(err.Error(), -1)
		}
		workflow := s.Workflow(workflowName)
		if workflow == nil {
			return cli.NewExitError("workflow not found", 404)
		}
		workflowID = workflow.ID
	default:
		return cli.NewExitError("must specify `<workflow name>`", -1)
	}

	res, err := retry.Workflow(ctx, client, workflowID)
	if err != nil {
		return cli.NewExitError(err.Error(), -1)
	}
	fmt.Println(res)
	return nil
}
