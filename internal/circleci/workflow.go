package circleci

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (c *Client) baseWorkflowURL() string {
	return fmt.Sprintf(
		"%s/api/v2/workflow",
		c.rootURL,
	)
}

type workflowJobListReponse struct {
	Items []*Job `json:"items"`
}

// https://circleci.com/docs/api/v2/#operation/listWorkflowJobs
func (c *Client) jobs(ctx context.Context, w *Workflow) ([]*Job, error) {
	url := c.baseWorkflowURL()

	url = fmt.Sprintf("%s/%s/job", url, w.ID)

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

	wjl := workflowJobListReponse{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&wjl); err != nil {
		return nil, err
	}

	return wjl.Items, nil
}

type retryResponse struct {
	Message string `json:"message"`
}

// RetryWorkflow will re-run the given workflow.  The endpoint documentation
// appears to allow specifying the job ids to presumably only run a single job
// in a workflow. However, this does not seem to be the behavior, instead it
// runs the full workflow even if a single job is specified.
//
// https://circleci.com/docs/api/v2/#operation/rerunWorkflow
func (c *Client) RetryWorkflow(ctx context.Context, workflowID string) (string, error) {
	url := c.baseWorkflowURL()

	url = fmt.Sprintf("%s/%s/rerun", url, workflowID)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := c.do(ctx, req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("response error: %d: %q", resp.StatusCode, body)
	}

	rr := retryResponse{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&rr); err != nil {
		return "", err
	}
	return rr.Message, nil
}
