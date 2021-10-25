package circleci

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// BuildSummaryWorkflow is used to Marshal the response from CircleCI
// when checking the status of a branch.
// Only fields used by cci are populated.
//
// https://circleci.com/docs/api/v1/?shell#recent-builds-for-a-single-project
type BuildSummaryWorkflow struct {
	JobName string `json:"job_name"`
	Name    string `json:"workflow_name"`
}

// BuildSummary is used to Marshal the response from CircleCI
// when checking the status of a branch.
// Only fields used by cci are populated.
//
// https://circleci.com/docs/api/v1/?shell#recent-builds-for-a-single-project
type BuildSummary struct {
	BuildNum  uint64               `json:"build_num"`
	Status    string               `json:"status"`
	Workflows BuildSummaryWorkflow `json:"workflows"`
}

// BuildSummaryResponse is used to Marshal the response from CircleCI
// when checking the status of a branch.
//
// https://circleci.com/docs/api/v1/?shell#recent-builds-for-a-single-project
type BuildSummaryResponse []*BuildSummary

// BuildSummary retrieves the summary of the latest build for the given branch.
//
// https://circleci.com/docs/api/v1/?shell#recent-builds-for-a-single-project
func (c *Client) BuildSummary(ctx context.Context, branch string) (BuildSummaryResponse, error) {
	url := c.baseURL()
	if branch != "" {
		url = fmt.Sprintf("%s/tree/%s", url, branch)
	}

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

	bsr := make(BuildSummaryResponse, 0)
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&bsr); err != nil {
		return nil, err
	}

	return bsr, nil
}

// BuildAction is used to Marshal the response from CircleCI
// when retreiving the output of a build.
// A BuildAction contains a URL that has to the full output
// of a single action.
//
// https://circleci.com/docs/api/v1/?shell#single-job
type BuildAction struct {
	OutputURL string `json:"output_url"`
	HasOutput bool   `json:"has_output"`
}

// BuildStep is used to Marshal the response from CircleCI
// when retreiving the output of a Build. A Build can have
// multiple steps.
//
// https://circleci.com/docs/api/v1/?shell#single-job
type BuildStep struct {
	Name    string         `json:"name"`
	Actions []*BuildAction `json:"actions"`
}

// BuildResponse is used to Marshal the response from CircleCI
// when retreiving the output of a Build.
//
// https://circleci.com/docs/api/v1/?shell#single-job
type BuildResponse struct {
	Steps []*BuildStep `json:"steps"`
}

// Build retrieves the build for the given build number.
//
// https://circleci.com/docs/api/v1/?shell#single-job
func (c *Client) Build(ctx context.Context, num uint64) (*BuildResponse, error) {
	url := fmt.Sprintf("%s/%d", c.baseURL(), num)

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

	br := BuildResponse{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&br); err != nil {
		return nil, err
	}

	return &br, nil
}

// actionOutputResponse is used to Marshal the response from CircleCI.
// This is the response from the url provided by the BuildAction.OutputURL.
type actionOutputResponse struct {
	Message string `json:"message"`
}

// BuildActionOutput is used to get the full output for a BuildAction.
// If the BuildAction has output, it will retrieve the full output
// from the OutputURL.
func (c *Client) BuildActionOutput(ctx context.Context, b *BuildAction) (string, error) {
	if !b.HasOutput {
		return "", nil
	}

	req, err := http.NewRequest("GET", b.OutputURL, nil)
	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("response error: %d: %q", resp.StatusCode, body)
	}

	aor := []actionOutputResponse{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&aor); err != nil {
		return "", err
	}
	var out string
	for _, a := range aor {
		out += a.Message
	}
	return out, nil
}
