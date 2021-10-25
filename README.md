# CCI

[![CircleCI](https://circleci.com/gh/tmessi/cci/tree/main.svg?style=svg)](https://circleci.com/gh/tmessi/cci/tree/main)

A cli for fetching the status and full output of CircleCI jobs.


## Install

```bash
go get -u github.com/tmessi/cci/cci
```

## Usage

`cci` is designed to have sane defaults
if run from within cloned git repository.
It will examine the current branch, and origin remote
to determine which project and branch to use for queries to CircleCI.
Thus to check the status of the current branch, just run:

```bash
cci
```

However,
it does require a CircleCI Token to authenticate the requests.
It is recommended to use the environment variable, `CIRCLECI_TOKEN`,
along with something like
[direnv](https://direnv.net/).

First create a
[Personal Access Token](https://circleci.com/docs/2.0/managing-api-tokens/#creating-a-personal-api-token)
and add it to your environment:

```bash
export CIRCLE_CI_TOKEN=<personal access token>
```

For more usage information and flags, see the help:

```bash
cci --help
```

## Autocompletion

For `bash` copy the `.bash_completion` file to `/etc/bash_completion.d/`
or to a location that is sourced from `~/.bashrc` or `~/.bash_profile`.
Or add the contents directly to `~/.bashrc` or `~/.bash_profile`.

For other shells, see the
[urfave/cli docs](https://github.com/urfave/cli/blob/master/docs/v2/manual.md).
For example for `zsh` follow the description
[here](https://github.com/urfave/cli/blob/master/docs/v2/manual.md#zsh-support).
