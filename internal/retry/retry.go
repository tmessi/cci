// Package retry is used to retrieve the full output of a build.
package retry

import (
	"context"
)

type client interface {
	RetryWorkflow(context.Context, string) (string, error)
}

// Workflow will retrun the given workflow.
func Workflow(ctx context.Context, c client, workflowID string) (string, error) {
	s, err := c.RetryWorkflow(ctx, workflowID)
	return s, err
}
