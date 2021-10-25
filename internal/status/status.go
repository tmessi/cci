// Package status is used to check the status of a branch.
package status

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/tmessi/cci/internal/circleci"
	"github.com/tmessi/cci/internal/status/internal/template"
)

// Known errors.
var (
	ErrNoBranch = errors.New("no branch provided")
)

// Status reports the status of a set of CI workflows for a branch.
type Status struct {
	Workflows []*Workflow
}

func (s *Status) String() string {
	return template.Render(s.Workflows)
}

// Workflow returns the Workflow with the given name.
// If no Workflow is found, it will return nil.
func (s *Status) Workflow(name string) *Workflow {
	for _, w := range s.Workflows {
		if w.Name == name {
			return w
		}
	}
	return nil
}

// Workflow is a group of Jobs with a name.
type Workflow struct {
	Name string
	Jobs []*Job
}

// Job returns the Job with the given name.
// If no Job is found, it will return nil.
func (w *Workflow) Job(name string) *Job {
	for _, j := range w.Jobs {
		if j.Name == name {
			return j
		}
	}
	return nil
}

// Job reports the status of a single CI job for branch.
type Job struct {
	Name         string
	Status       string
	BuildNum     uint64
	WorkflowName string
}

func (j *Job) String() string {
	return fmt.Sprintf("%d: %s %s %s", j.BuildNum, j.WorkflowName, j.Name, j.Status)
}

type client interface {
	BuildSummary(context.Context, string) (circleci.BuildSummaryResponse, error)
}

// Check queries CircleCI for the status of a set of Jobs for the given Project and branch.
func Check(ctx context.Context, c client, branch string) (*Status, error) {
	if branch == "" {
		return nil, ErrNoBranch
	}

	bsr, err := c.BuildSummary(ctx, branch)
	if err != nil {
		return nil, err
	}

	jobsByWorkflowName := make(map[string][]*Job)
	for _, jr := range bsr {
		jobs, ok := jobsByWorkflowName[jr.Workflows.Name]
		if !ok {
			jobs = make([]*Job, 0)
		}
		jobs = append(jobs, &Job{
			Name:         jr.Workflows.JobName,
			Status:       jr.Status,
			BuildNum:     jr.BuildNum,
			WorkflowName: jr.Workflows.Name,
		})
		jobsByWorkflowName[jr.Workflows.Name] = jobs
	}

	workflows := make([]*Workflow, 0, len(jobsByWorkflowName))
	for name, jobs := range jobsByWorkflowName {
		sort.SliceStable(jobs, func(i, j int) bool { return jobs[i].Name < jobs[j].Name })

		workflows = append(workflows, &Workflow{
			Name: name,
			Jobs: jobs,
		})
	}
	sort.SliceStable(workflows, func(i, j int) bool { return workflows[i].Name < workflows[j].Name })

	return &Status{Workflows: workflows}, nil
}
