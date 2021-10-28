package circleci

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Job provides a summary of a Job. A workflow contains one or more Jobs.
type Job struct {
	ID     string `json:"id"`
	Number uint64 `json:"job_number"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

// Workflow provides a summary of a Workflow. A pipeline is made up of one or
// more workflows. Each workflow has one or more Jobs.
type Workflow struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Jobs   []*Job `json:"-"`
}

// Pipeline provides a summary of a single pipeline execution.
type Pipeline struct {
	ID        string      `json:"id"`
	Workflows []*Workflow `json:"-"`
}

func (c *Client) basePipelineListURL() string {
	return fmt.Sprintf(
		"%s/api/v2/project/%s/%s/%s/pipeline",
		c.rootURL,
		c.project.VCSType,
		c.project.Organization,
		c.project.Name,
	)
}

func (c *Client) basePipelineURL() string {
	return fmt.Sprintf(
		"%s/api/v2/pipeline",
		c.rootURL,
	)
}

type pipelineListResponse struct {
	Items []*Pipeline `json:"items"`
}

// https://circleci.com/docs/api/v2/#operation/listPipelinesForProject
func (c *Client) mostRecentPipeline(ctx context.Context, branch string) (*Pipeline, error) {
	url := c.basePipelineListURL()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if branch != "" {
		q := req.URL.Query()
		q.Add("branch", branch)
		req.URL.RawQuery = q.Encode()
	}

	resp, err := c.do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("response error: %d: %q", resp.StatusCode, body)
	}

	plr := pipelineListResponse{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&plr); err != nil {
		return nil, err
	}

	if len(plr.Items) <= 0 {
		return nil, fmt.Errorf("no pipelines for branch: %s", branch)
	}

	return plr.Items[0], nil
}

type pipelineWorkflowListReponse struct {
	Items []*Workflow `json:"items"`
}

// https://circleci.com/docs/api/v2/#operation/listWorkflowsByPipelineId
func (c *Client) workflows(ctx context.Context, p *Pipeline) ([]*Workflow, error) {
	url := c.basePipelineURL()

	url = fmt.Sprintf("%s/%s/workflow", url, p.ID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("response error: %d: %q", resp.StatusCode, body)
	}

	pwlr := pipelineWorkflowListReponse{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&pwlr); err != nil {
		return nil, err
	}

	return pwlr.Items, nil
}

// PipelineSummary returns a summary of the most recent Pipeline execution for the given branch.
func (c *Client) PipelineSummary(ctx context.Context, branch string) (*Pipeline, error) {
	p, err := c.mostRecentPipeline(ctx, branch)
	if err != nil {
		return nil, err
	}

	p.Workflows, err = c.workflows(ctx, p)
	if err != nil {
		return nil, err
	}

	for _, w := range p.Workflows {
		w.Jobs, err = c.jobs(ctx, w)
		if err != nil {
			return nil, err
		}
	}

	return p, nil
}
