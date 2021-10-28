// Package status is used to check the status of a branch.
package status

import (
	"context"
	"errors"

	"github.com/tmessi/cci/internal/circleci"
	"github.com/tmessi/cci/internal/status/internal/template"
)

// Known errors.
var (
	ErrNoBranch = errors.New("no branch provided")
)

// Status reports the status of a set of CI workflows for a branch.
type Status struct {
	Pipelines []*circleci.Pipeline
}

func (s *Status) String() string {
	return template.Render(s)
}

// Workflow returns the Workflow with the given name.
// If no Workflow is found, it will return nil.
func (s *Status) Workflow(name string) *circleci.Workflow {
	if len(s.Pipelines) <= 0 {
		return nil
	}

	for _, w := range s.Pipelines[0].Workflows {
		if w.Name == name {
			return w
		}
	}
	return nil
}

// Job returns the Job with the given name.
// If no Job is found, it will return nil.
func (s *Status) Job(workflowName, name string) *circleci.Job {
	w := s.Workflow(workflowName)
	if w != nil {
		for _, j := range w.Jobs {
			if j.Name == name {
				return j
			}
		}
	}
	return nil
}

type client interface {
	Pipelines(context.Context, string, uint64) ([]*circleci.Pipeline, error)
}

// Check queries CircleCI for the status of a set of Jobs for the given Project and branch.
func Check(ctx context.Context, c client, branch string, limit uint64) (*Status, error) {
	if branch == "" {
		return nil, ErrNoBranch
	}

	if limit == 0 {
		limit = 1
	}

	p, err := c.Pipelines(ctx, branch, limit)
	if err != nil {
		return nil, err
	}

	return &Status{Pipelines: p}, nil
}
