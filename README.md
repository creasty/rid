rid (run-in-docker)
===================

[![Build Status](https://travis-ci.org/creasty/rid.svg?branch=master)](https://travis-ci.org/creasty/rid) [![codecov](https://codecov.io/gh/creasty/rid/branch/master/graph/badge.svg)](https://codecov.io/gh/creasty/rid)

Run commands in container as if were native. Stress-free dockerized development environment finally arrived.


What is rid?
------------

With a rid directory on a root directory of project, any command prefixed by `rid` is executed within a Docker container.

```hcl
$ ls ./rid
config.yml
docker-compose.yml
Dockerfile
```

Lo and behold, even if your environment is absolutely clean and there's only `docker` and `docker-compose` are installed,
all you have to do is the followings so as to get started with a new project from scratch.


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

### MacOS

```hcl
$ brew install creasty/tools/rid
```

### Other

Download from here:  
https://github.com/creasty/rid/releases
