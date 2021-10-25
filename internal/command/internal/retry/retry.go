// Package retry provides the retry subcommand.
package retry

import (
	"fmt"
	"strconv"

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
	ArgsUsage:    "<build number> | <workflow name> <job name>",
	Aliases:      []string{"r"},
	Usage:        "Retry a build",
	BashComplete: complete.Build,
	Action:       action,
}

func action(c *cli.Context) error {
	ctx, cancel := signal.InitContext()
	defer cancel()

	client, err := global.Client(c)
	if err != nil {
		return cli.NewExitError(err.Error(), -1)
	}

	var n uint64

	switch c.NArg() {
	case 2:
		workflowName := c.Args().Get(0)
		jobName := c.Args().Get(1)

		s, err := status.Check(ctx, client, c.String("branch"))
		if err != nil {
			return cli.NewExitError(err.Error(), -1)
		}
		workflow := s.Workflow(workflowName)
		if workflow == nil {
			return cli.NewExitError("workflow not found", 404)
		}
		job := workflow.Job(jobName)
		if job == nil {
			return cli.NewExitError("job not found", 404)
		}
		n = job.BuildNum
	case 1:
		strconv.Atoi(c.Args().Get(0))
		buildNum, err := strconv.Atoi(c.Args().Get(0))
		if err != nil {
			return cli.NewExitError("build number must be an int", -1)
		}
		n = uint64(buildNum)
	default:
		return cli.NewExitError("must specify `<build number>` or `<workflow name> <job name>`", -1)
	}

	b, err := retry.Build(ctx, client, n)
	if err != nil {
		return cli.NewExitError(err.Error(), -1)
	}
	fmt.Println(b)
	return nil
}
