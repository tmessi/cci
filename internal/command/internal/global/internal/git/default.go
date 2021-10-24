// Package git provides default values for global flags
// by looking for git information. This allows for providing
// sane defaults if cci is run from within a git repository.
// It will allow cci to check CircleCI for the current branch of
// the repo in the current working directory.
package git

import (
	"os"
	"strings"

	"github.com/go-git/go-git/v5"

	giturls "github.com/whilp/git-urls"
)

type defaults struct {
	Type         string
	Organization string
	Repository   string
	Branch       string
}

// Defaults provide default values for some global flags.
// It is populated via calls to git, allowing for sane defaults
// if run from within a git repository.
var Defaults = &defaults{}

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		return
	}

	repo, err := git.PlainOpenWithOptions(
		cwd,
		&git.PlainOpenOptions{DetectDotGit: true},
	)
	if err != nil {
		return
	}

	remote, err := repo.Remote("origin")
	if err != nil {
		return
	}

	remoteURL := remote.Config().URLs[0]
	u, err := giturls.Parse(remoteURL)
	if err != nil {
		return
	}

	switch u.Host {
	case "bitbucket.org":
		Defaults.Type = "bitbucket"
	case "github.com":
		fallthrough
	default:
		Defaults.Type = "github"
	}

	parts := strings.SplitN(u.Path, "/", 2)
	if len(parts) != 2 {
		return
	}
	Defaults.Organization = parts[0]
	Defaults.Repository = strings.TrimSuffix(parts[1], ".git")

	ref, err := repo.Head()
	if err != nil {
		return
	}

	if ref.Name().IsBranch() {
		Defaults.Branch = ref.Name().Short()
	}
}
