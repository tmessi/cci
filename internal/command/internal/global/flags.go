// Package global provides flags for all subcommands.
package global

import (
	"errors"
	"net/http"

	"github.com/tmessi/cci/internal/circleci"
	"github.com/tmessi/cci/internal/command/internal/global/internal/git"
	"gopkg.in/urfave/cli.v1"
)

// Flags are the flags that apply to subcommands.
var Flags = []cli.Flag{
	cli.StringFlag{
		Name:   "circleci-token,token",
		Usage:  "Token used for authenticating to the CircleCI API.",
		EnvVar: "CIRCLE_CI_TOKEN",
	},
	cli.StringFlag{
		Name:   "vcs-type",
		Usage:  "The vcs type for the project.",
		EnvVar: "PROJECT_VCS_TYPE",
		Value:  git.Defaults.Type,
	},
	cli.StringFlag{
		Name:   "org",
		Usage:  "The organization for the project.",
		EnvVar: "PROJECT_ORG",
		Value:  git.Defaults.Organization,
	},
	cli.StringFlag{
		Name:   "project",
		Usage:  "The project.",
		EnvVar: "PROJECT",
		Value:  git.Defaults.Repository,
	},
	cli.StringFlag{
		Name:   "branch,b",
		Usage:  "The branch to check",
		EnvVar: "CCI_BRANCH",
		Value:  git.Defaults.Branch,
	},
}

// Errors for invalid flag values.
var (
	ErrNoProject = errors.New("no project specified")
	ErrNoOrg     = errors.New("no organization specified")
	ErrNoVCSType = errors.New("no vsc-type specified")
	ErrNoToken   = errors.New("no circleci tokoen specified")
)

// Client creates a circleci.Client from the global cli flags.
func Client(c *cli.Context) (*circleci.Client, error) {
	project := c.GlobalString("project")
	org := c.GlobalString("org")
	vcsType := c.GlobalString("vcs-type")
	token := c.GlobalString("token")

	if project == "" {
		return nil, ErrNoProject
	}

	if org == "" {
		return nil, ErrNoOrg
	}

	if vcsType == "" {
		return nil, ErrNoVCSType
	}

	if token == "" {
		return nil, ErrNoToken
	}

	return circleci.New(
		&http.Client{},
		&circleci.Project{
			Name:         project,
			Organization: org,
			VCSType:      vcsType,
		},
		token,
	), nil
}
