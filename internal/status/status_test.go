package status_test

import (
	"context"
	"errors"
	"testing"

	"github.com/tmessi/cci/internal/circleci"
	"github.com/tmessi/cci/internal/status"
)

type testClient struct {
	bsr circleci.BuildSummaryResponse
	err error
}

func (c *testClient) BuildSummary(_ context.Context, _ string) (circleci.BuildSummaryResponse, error) {
	return c.bsr, c.err
}

func TestCheck(t *testing.T) {
	tests := []struct {
		name           string
		bsr            circleci.BuildSummaryResponse
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
			circleci.BuildSummaryResponse{
				{
					BuildNum: 3,
					Status:   "success",
					Workflows: circleci.BuildSummaryWorkflow{
						JobName: "unit",
						Name:    "tests",
					},
				},
				{
					BuildNum: 4,
					Status:   "success",
					Workflows: circleci.BuildSummaryWorkflow{
						JobName: "gofmt",
						Name:    "lint",
					},
				},
				{
					BuildNum: 5,
					Status:   "success",
					Workflows: circleci.BuildSummaryWorkflow{
						JobName: "govet",
						Name:    "lint",
					},
				},
			},
			nil,
			"branch",
			&status.Status{
				Workflows: []*status.Workflow{
					{
						Name: "lint",
						Jobs: []*status.Job{
							{
								"gofmt",
								"success",
								4,
								"lint",
							},
							{
								"govet",
								"success",
								5,
								"lint",
							},
						},
					},
					{
						Name: "tests",
						Jobs: []*status.Job{
							{
								"unit",
								"success",
								3,
								"tests",
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
			client := &testClient{tt.bsr, tt.err}

			ctx := context.Background()
			s, err := status.Check(ctx, client, tt.branch)

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

			if len(s.Workflows) != len(tt.expectedStatus.Workflows) {
				t.Fatalf("got %d workflows, wanted %d workflows", len(s.Workflows), len(tt.expectedStatus.Workflows))
			}

			for i, w := range s.Workflows {
				wantWorkflow := tt.expectedStatus.Workflows[i]

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

					if got.BuildNum != want.BuildNum {
						t.Errorf("BuildNum: got %d, wanted %d", got.BuildNum, want.BuildNum)
					}

					if got.WorkflowName != want.WorkflowName {
						t.Errorf("WorkflowName: got %q, wanted %q", got.WorkflowName, want.WorkflowName)
					}
				}
			}
		})
	}
}
