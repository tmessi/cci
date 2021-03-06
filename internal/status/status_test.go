package status_test

import (
	"context"
	"errors"
	"testing"

	"github.com/tmessi/cci/internal/circleci"
	"github.com/tmessi/cci/internal/status"
)

type testClient struct {
	pipelines []*circleci.Pipeline
	err       error
}

func (c *testClient) Pipelines(_ context.Context, _ string, _ uint64) ([]*circleci.Pipeline, error) {
	return c.pipelines, c.err
}

func TestCheck(t *testing.T) {
	tests := []struct {
		name           string
		pipelines      []*circleci.Pipeline
		err            error
		branch         string
		expectedStatus *status.Status
		expectedError  error
	}{
		{
			"Error",
			nil,
			errors.New("error"),
			"branch",
			nil,
			errors.New("error"),
		},
		{
			"NoBranch",
			nil,
			errors.New("error"),
			"",
			nil,
			status.ErrNoBranch,
		},
		{
			"MultipleJobs",
			[]*circleci.Pipeline{
				{
					Workflows: []*circleci.Workflow{
						{
							ID:     "11111111-1111-1111-1111-111111111111",
							Name:   "tests",
							Status: "success",
							Jobs: []*circleci.Job{
								{
									ID:     "11111111-1111-1111-1111-111111111112",
									Name:   "unit",
									Number: 1,
									Status: "success",
								},
							},
						},
						{
							ID:     "11111111-1111-1111-2222-111111111111",
							Name:   "lint",
							Status: "success",
							Jobs: []*circleci.Job{
								{
									ID:     "11111111-1111-1111-2222-111111111112",
									Name:   "gofmt",
									Number: 2,
									Status: "success",
								},
								{
									ID:     "11111111-1111-1111-2222-111111111113",
									Name:   "govet",
									Number: 3,
									Status: "success",
								},
							},
						},
					},
				},
			},
			nil,
			"branch",
			&status.Status{
				Pipelines: []*circleci.Pipeline{
					{
						Workflows: []*circleci.Workflow{
							{
								ID:     "11111111-1111-1111-1111-111111111111",
								Name:   "tests",
								Status: "success",
								Jobs: []*circleci.Job{
									{
										ID:     "11111111-1111-1111-1111-111111111112",
										Name:   "unit",
										Number: 1,
										Status: "success",
									},
								},
							},
							{
								ID:     "11111111-1111-1111-2222-111111111111",
								Name:   "lint",
								Status: "success",
								Jobs: []*circleci.Job{
									{
										ID:     "11111111-1111-1111-2222-111111111112",
										Name:   "gofmt",
										Number: 2,
										Status: "success",
									},
									{
										ID:     "11111111-1111-1111-2222-111111111113",
										Name:   "govet",
										Number: 3,
										Status: "success",
									},
								},
							},
						},
					},
				},
			},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &testClient{tt.pipelines, tt.err}

			ctx := context.Background()
			s, err := status.Check(ctx, client, tt.branch, uint64(len(tt.pipelines)))

			if tt.expectedError != nil {
				if err == nil {
					t.Fatalf("expected error, but did not get one")
				}
				if tt.expectedError.Error() != err.Error() {
					t.Errorf("got %q, wanted %q", err.Error(), tt.expectedError.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("got error, but did not expect one")
			}

			if len(s.Pipelines) != len(tt.expectedStatus.Pipelines) {
				t.Fatalf("got %d pipelines, wanted %d pipelines", len(s.Pipelines), len(tt.expectedStatus.Pipelines))
			}

			for i, p := range s.Pipelines {
				wantPipeline := tt.expectedStatus.Pipelines[i]

				if len(p.Workflows) != len(wantPipeline.Workflows) {
					t.Fatalf("got %d workflows, wanted %d workflows", len(p.Workflows), len(wantPipeline.Workflows))
				}

				for i, w := range p.Workflows {
					wantWorkflow := wantPipeline.Workflows[i]

					if len(w.Jobs) != len(wantWorkflow.Jobs) {
						t.Fatalf("got %d jobs, wanted %d jobs", len(w.Jobs), len(wantWorkflow.Jobs))
					}

					if w.Name != wantWorkflow.Name {
						t.Errorf("Name: got %q, wanted %q", w.Name, wantWorkflow.Name)
					}

					for j, got := range w.Jobs {
						want := wantWorkflow.Jobs[j]

						if got.Name != want.Name {
							t.Errorf("Name: got %q, wanted %q", got.Name, want.Name)
						}

						if got.Status != want.Status {
							t.Errorf("Status: got %q, wanted %q", got.Status, want.Status)
						}

						if got.Number != want.Number {
							t.Errorf("Number: got %d, wanted %d", got.Number, want.Number)
						}
					}
				}
			}
		})
	}
}
