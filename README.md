![](https://cloud.githubusercontent.com/assets/1695538/24829859/002571a4-1cb5-11e7-8e1b-c9b04d171828.png)


rid (run-in-docker)
===================

[![Build Status](https://travis-ci.org/creasty/rid.svg?branch=master)](https://travis-ci.org/creasty/rid)
[![GitHub release](https://img.shields.io/github/release/creasty/rid.svg)](https://github.com/creasty/rid/releases)
[![License](https://img.shields.io/github/license/creasty/rid.svg)](./LICENSE)

Run commands in container as if were native. Stress-free dockerized development environment finally arrived.


What is rid?
------------

With a `rid/` directory at the root of a project, any command prefixed by `rid` is executed within a Docker container.

```hcl
$ ls ./rid
config.yml
docker-compose.yml
Dockerfile
```

That is to say, even if your environment is absolutely clean and you have nothing but `docker`, `docker-compose` and `rid`, getting started with a new Rails project from scratch has never been easier.

```hcl
# install dependencies and setup a database
$ rid cp .env{.sample,}
$ rid bundle install --path vendor/bundle
$ rid rake db:create
$ rid rake db:schema:load

# start a server
$ rid rails s
```


Installation
------------

First, install [Docker](https://docs.docker.com/engine/installation/) and [Docker Compose](https://docs.docker.com/compose/install/).
The easiest way to do this on macOS is by installing [Docker for Mac](https://docs.docker.com/docker-for-mac/).

### macOS

You can install `rid` via Homebrew:

```hcl
$ brew install creasty/tools/rid
```

### Linux

Download a binary from here: https://github.com/creasty/rid/releases

### Windows

Not supported yet


Usage
-----

`rid` is a project contextual tool, meaning that it's aware of working directory and automatically finds the root directory of a project by locating a configuration file.

Typical `rid` directory looks like this:

```hcl
rid/                   # rid directory at the root of a project (e.g., same level as `.git`'s)
  libexec/             # custom sub-commands for rid
  config.yml           # configuration file for rid
  docker-compose.yml   # docker-compose manifest
  Dockerfile           # dockerfile
```

Note that `rid/config.yml` and `rid/docker-compose.yml` are regardlessly required for `rid` to work with.

### Config file

Configurable parameters of `rid/config.yml` are the following.

```go
type Config struct {
	// ProjectName is used for `docker-compose` in order to distinguish projects in other locations
	ProjectName string `json:"project_name" valid:"required"`

	// MainService is a service name in `docker-compose.yml`, in which container commands given to rid are executed
	// Default is "app"
	MainService string `json:"main_service"`
}
```

### Custom commands

Executables in `rid/libexec/` can be run as a sub command.

```hcl
rid/libexec/
  foo           # `rid foo` -- this is executed in a container
  rid-bar       # `rid bar` -- name starts with `rid-` is executed on a host computer
  rid-bar.txt   # optionally, placing `.txt` file that shares the common basename enables "help" functionality
```

Help file should have a title in the first line:

```
Show greeting message

Usage:
    rid bar NAME
```

The title (first line) appears on the help of `rid`.

```hcl
$ rid
Execute commands via docker-compose

Usage:
    rid COMMAND [args...]
    rid COMMAND -h|--help
    rid [options]

Options:
    -h, --help     Show this
    -v, --version  Show rid version
        --debug    Debug context and configuration

Commands:
    compose  # Execute docker-compose
    foo
    bar      # Show greeting message
```

And `rid COMMAND -h` prints the full contents.

```
$ rid bar -h
Show greeting message

Usage:
    rid bar NAME
```


Development
-----------

Surprise surprise, `rid` itself is developed by `rid`!

```hcl
$ rid glide install  # install dependencies
$ rid make test      # run lint and tests
$ rid make           # compile for darwin/amd64
$ ./bin/rid -v       # execute a new binary
```
