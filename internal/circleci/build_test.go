package circleci_test

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/tmessi/cci/internal/circleci"
)

func TestBuild(t *testing.T) {
	tests := []struct {
		name     string
		project  *circleci.Project
		token    string
		num      uint64
		err      error
		expected *circleci.BuildResponse
	}{
		{
			"MultipleSteps",
			&circleci.Project{
				Name:         "cci",
				Organization: "tmessi",
				VCSType:      "github",
			},
			"valid-token",
			1,
			nil,
			&circleci.BuildResponse{
				Steps: []*circleci.BuildStep{
					{
						Name: "Spin up environment",
						Actions: []*circleci.BuildAction{
							{
								OutputURL: "https://circle-production-action-output.s3.amazonaws.com/b4a9ff1942198703d5db5716-61745bc6bc80122676db990e-0-0",
								HasOutput: true,
							},
						},
					},
					{
						Name: "Preparing environment variables",
						Actions: []*circleci.BuildAction{
							{
								OutputURL: "https://circle-production-action-output.s3.amazonaws.com/6da08a51414f17bc06db5716-61745bc6bc80122676db990e-99-0",
								HasOutput: true,
							},
						},
					},
					{
						Name: "Checkout code",
						Actions: []*circleci.BuildAction{
							{
								OutputURL: "https://circle-production-action-output.s3.amazonaws.com/7da08a51414f17bc06db5716-61745bc6bc80122676db990e-101-0",
								HasOutput: true,
							},
						},
					},
					{
						Name: "Restoring cache",
						Actions: []*circleci.BuildAction{
							{
								OutputURL: "https://circle-production-action-output.s3.amazonaws.com/05a9ff194219870316db5716-61745bc6bc80122676db990e-102-0",
								HasOutput: true,
							},
						},
					},
					{
						Name: "Install Dependencies",
						Actions: []*circleci.BuildAction{
							{
								OutputURL: "",
								HasOutput: false,
							},
						},
					},
					{
						Name: "Saving cache",
						Actions: []*circleci.BuildAction{
							{
								OutputURL: "https://circle-production-action-output.s3.amazonaws.com/95a9ff194219870376db5716-61745bc6bc80122676db990e-104-0",
								HasOutput: true,
							},
						},
					},
					{
						Name: "Install golangci-lint",
						Actions: []*circleci.BuildAction{
							{
								OutputURL: "https://circle-production-action-output.s3.amazonaws.com/a5a9ff194219870376db5716-61745bc6bc80122676db990e-105-0",
								HasOutput: true,
							},
						},
					},
					{
						Name: "Run lint",
						Actions: []*circleci.BuildAction{
							{
								OutputURL: "",
								HasOutput: false,
							},
						},
					},
				},
			},
		},
		{
			"BadToken",
			&circleci.Project{
				Name:         "cci",
				Organization: "tmessi",
				VCSType:      "github",
			},
			"invalid-token",
			1,
			errors.New(`response error: 404: "{\n  \"message\" : \"Project not found\"\n}\n"`),
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				code, err := ioutil.ReadFile(fmt.Sprintf("testdata/%s.code", t.Name()))
				if err != nil {
					t.Fatalf("test not configured correctly: %s", err.Error())
				}
				statusCode, err := strconv.Atoi(strings.TrimSpace(string(code)))
				if err != nil {
					t.Fatalf("test not configured correctly: %s", err.Error())
				}
				res, err := ioutil.ReadFile(fmt.Sprintf("testdata/%s.json", t.Name()))
				if err != nil {
					t.Fatalf("test not configured correctly: %s", err.Error())
				}

				w.WriteHeader(statusCode)
				w.Write(res)
			}))
			defer ts.Close()

			tc := ts.Client()

			client := circleci.New(tc, ts.URL, tt.project, tt.token)

			ctx := context.Background()
			br, err := client.Build(ctx, tt.num)

			if tt.err != nil {
				if err == nil {
					t.Fatalf("did not get error but expected: %s", tt.err.Error())
				}
				if err.Error() != tt.err.Error() {
					t.Errorf("got %q, wanted %q", err.Error(), tt.err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("err: %s", err.Error())
			}

			if len(br.Steps) != len(tt.expected.Steps) {
				t.Fatalf("Steps: got %d, wanted %d", len(br.Steps), len(tt.expected.Steps))
			}

			for i, gotStep := range br.Steps {
				wantStep := tt.expected.Steps[i]

				if gotStep.Name != wantStep.Name {
					t.Errorf("Name: got %q, want %q", gotStep.Name, wantStep.Name)
				}

				if len(gotStep.Actions) != len(wantStep.Actions) {
					t.Fatalf("Actions: got %d, wanted %d", len(gotStep.Actions), len(wantStep.Actions))
				}

				for j, gotAction := range gotStep.Actions {
					wantAction := wantStep.Actions[j]

					if gotAction.HasOutput != wantAction.HasOutput {
						t.Errorf("HasOutput: got %t, want %t", gotAction.HasOutput, wantAction.HasOutput)
					}

					if gotAction.OutputURL != wantAction.OutputURL {
						t.Errorf("OutputURL: got %q, want %q", gotAction.OutputURL, wantAction.OutputURL)
					}

				}
			}
		})
	}
}

func TestBuildActionOutput(t *testing.T) {
	tests := []struct {
		name     string
		project  *circleci.Project
		token    string
		action   *circleci.BuildAction
		err      error
		expected string
	}{
		{
			"HasOutputTrue",
			&circleci.Project{
				Name:         "cci",
				Organization: "tmessi",
				VCSType:      "github",
			},
			"valid-token",
			&circleci.BuildAction{
				HasOutput: true,
				OutputURL: "/out.json",
			},
			nil,
			"Found a cache from build 13 at go-mod-v4-e9cX1Vj+6JCCGWLDOuWSzyiSzF3fE9ayVDQOCbRkWAM=\nSize: 226 MiB\nCached paths:\n  * /go/pkg/mod\n\nDownloading cache archive...\nValidating cache...\n\nUnarchiving cache...\n",
		},
		{
			"HasOutputFalse",
			&circleci.Project{
				Name:         "cci",
				Organization: "tmessi",
				VCSType:      "github",
			},
			"valid-token",
			&circleci.BuildAction{
				HasOutput: false,
				OutputURL: "",
			},
			nil,
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				code, err := ioutil.ReadFile(fmt.Sprintf("testdata/%s.code", t.Name()))
				if err != nil {
					t.Fatalf("test not configured correctly: %s", err.Error())
				}
				statusCode, err := strconv.Atoi(strings.TrimSpace(string(code)))
				if err != nil {
					t.Fatalf("test not configured correctly: %s", err.Error())
				}
				res, err := ioutil.ReadFile(fmt.Sprintf("testdata/%s.json", t.Name()))
				if err != nil {
					t.Fatalf("test not configured correctly: %s", err.Error())
				}

				w.WriteHeader(statusCode)
				w.Write(res)
			}))
			defer ts.Close()

			tc := ts.Client()

			client := circleci.New(tc, ts.URL, tt.project, tt.token)

			ctx := context.Background()
			tt.action.OutputURL = ts.URL + tt.action.OutputURL
			out, err := client.BuildActionOutput(ctx, tt.action)

			if tt.err != nil {
				if err == nil {
					t.Fatalf("did not get error but expected: %s", tt.err.Error())
				}
				if err.Error() != tt.err.Error() {
					t.Errorf("got %q, wanted %q", err.Error(), tt.err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("err: %s", err.Error())
			}

			if out != tt.expected {
				t.Errorf("got %q, want %q", out, tt.expected)
			}
		})
	}
}
