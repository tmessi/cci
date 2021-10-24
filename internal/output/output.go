// Package output is used to retrieve the full output of a build.
package output

import (
	"context"
	"fmt"

	"github.com/tmessi/cci/internal/circleci"
	"github.com/tmessi/cci/internal/output/internal/template"
)

// Step is a stage in a build. It has a name and the output
// of the step.
type Step struct {
	Name   string
	Output string
}

func (s *Step) String() string {
	return fmt.Sprintf("%s:\n%s", s.Name, s.Output)
}

// Build is a set of Steps.
type Build struct {
	Steps []*Step
}

func (b *Build) String() string {
	return template.Render(b)
}

type client interface {
	Build(context.Context, uint64) (*circleci.BuildResponse, error)
	BuildActionOutput(context.Context, *circleci.BuildAction) (string, error)
}

// GetBuild retrieves the output of given build number.
func GetBuild(ctx context.Context, c client, buildNum uint64) (*Build, error) {
	br, err := c.Build(ctx, buildNum)
	if err != nil {
		return nil, err
	}

	steps := make([]*Step, 0, len(br.Steps))
	for _, step := range br.Steps {
		for _, action := range step.Actions {
			if action.HasOutput {
				o, err := c.BuildActionOutput(ctx, action)
				if err != nil {
					return nil, err
				}
				steps = append(steps, &Step{
					Name:   step.Name,
					Output: o,
				})
			}
		}
	}
	return &Build{Steps: steps}, nil
}
