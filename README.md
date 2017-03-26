rid (run-in-docker)
===================

[![Build Status](https://travis-ci.org/creasty/rid.svg?branch=master)](https://travis-ci.org/creasty/rid)
[![codecov](https://codecov.io/gh/creasty/rid/branch/master/graph/badge.svg)](https://codecov.io/gh/creasty/rid)
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

That is to say, even if your environment is absolutely clean and all you have is only `docker` and `docker-compose`, getting started with a new Rails project from scratch has never been easier.

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

First, install [Docker](https://docs.docker.com/engine/installation/) and [Docker Compose](https://docs.docker.com/compose/install/). The easiest way to do this on MacOS is by installing [Docker for Mac](https://docs.docker.com/docker-for-mac/).

### MacOS

```hcl
$ brew install creasty/tools/rid
```

### Linux

Download it from here: https://github.com/creasty/rid/releases

### Windows

Not supported yet


Development
-----------

`rid` itself is also developed by `rid`.

```hcl
$ rid glide install  # install dependencies
$ rid make test      # run lint and tests
$ rid make           # compile for darwin/amd64
$ ./bin/rid -v
```
