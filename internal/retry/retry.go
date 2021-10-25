// Package retry is used to retrieve the full output of a build.
package retry

import (
	"context"
	"fmt"

	"github.com/tmessi/cci/internal/circleci"
)

// Summary provides details about the new build.
type Summary struct {
	BuildNum uint64
}

func (s *Summary) String() string {
	return fmt.Sprintf("%d\n", s.BuildNum)
}

type client interface {
	RetryBuild(context.Context, uint64) (*circleci.BuildSummary, error)
}

// Build will execute a build again.
func Build(ctx context.Context, c client, buildNum uint64) (*Summary, error) {
	br, err := c.RetryBuild(ctx, buildNum)
	if err != nil {
		return nil, err
	}
	return &Summary{
		BuildNum: br.BuildNum,
	}, nil
}
