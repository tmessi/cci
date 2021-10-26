// Package global provides flags for all subcommands.
package global

import (
	"errors"
	"net/http"

	"github.com/tmessi/cci/internal/circleci"
	"github.com/tmessi/cci/internal/command/internal/global/internal/git"
	"github.com/urfave/cli/v2"
)

// Flags are the flags that apply to subcommands.
var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:    "token",
		Aliases: []string{"circleci-token"},
		Usage:   "Token used for authenticating to the CircleCI API.",
		EnvVars: []string{"CIRCLE_CI_TOKEN"},
	},
	&cli.StringFlag{
		Name:    "vcs-type",
		Usage:   "The vcs type for the project.",
		EnvVars: []string{"PROJECT_VCS_TYPE"},
		Value:   git.Defaults.Type,
	},
	&cli.StringFlag{
		Name:    "org",
		Usage:   "The organization for the project.",
		EnvVars: []string{"PROJECT_ORG"},
		Value:   git.Defaults.Organization,
	},
	&cli.StringFlag{
		Name:    "project",
		Usage:   "The project.",
		EnvVars: []string{"PROJECT"},
		Value:   git.Defaults.Repository,
	},
	&cli.StringFlag{
		Name:    "branch",
		Aliases: []string{"b"},
		Usage:   "The branch to check",
		EnvVars: []string{"CCI_BRANCH"},
		Value:   git.Defaults.Branch,
	},
	&cli.StringFlag{
		Name:    "url",
		Aliases: []string{"root-url"},
		Usage:   "The root url for the circleci api",
		EnvVars: []string{"CIRCLE_CI_URL"},
		Value:   "https://circleci.com",
	},
}

// Errors for invalid flag values.
var (
	ErrNoProject = errors.New("no project specified")
	ErrNoOrg     = errors.New("no organization specified")
	ErrNoVCSType = errors.New("no vcs-type specified")
	ErrNoToken   = errors.New("no circleci token specified")
)

// Client creates a circleci.Client from the global cli flags.
func Client(c *cli.Context) (*circleci.Client, error) {
	project := c.String("project")
	org := c.String("org")
	vcsType := c.String("vcs-type")
	token := c.String("token")

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
		c.String("url"),
		&circleci.Project{
			Name:         project,
			Organization: org,
			VCSType:      vcsType,
		},
		token,
	), nil
}
