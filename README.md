Go monitoring tools
===================

*Augitne.com* specific tools for dealing with metrics and logs transformation for `ElasticSearch`/`Kibana`.

Using `Golang` because `Logstash` takes too much resources and `golang` can download/parse/publish logs in parallel.

Building with docker
--------------------

 * Install [docker](https://docs.docker.com/engine/installation/)
 * Install [glide](https://github.com/Masterminds/glide)
 * Downlaod dependencies: `glide install`
 * Run `scripts/build.sh`
 * Check results in `bin` folder
 
Using docker-compose
--------------------

```
version: "2"

services:
  auginte.dev.gologs:
    image: golang:1.7.0
    volumes:
      - ./:/go/src/github.com/Auginte/go-monitoring/
      - ./bin:/go/bin
    command: go install github.com/Auginte/go-monitoring/...
```

Assuming `docker-compose.yml` file is in current directory (otherwise updates `volumes` section)
 
Developing with local go
------------------------

 * [Install go 1.7](https://golang.org/doc/install)
 * clone this project into `$GOPATH/src/github.com/Auginte/go-monitoring/`
 * [Run from IDE](https://plugins.jetbrains.com/plugin/5047)

Tests?
------

Currently everything is tested manually: on Ubuntu and Amazon AMI.
Only small part of tests are used only to check, if it compiles.
[![Build Status](https://travis-ci.org/Auginte/go-monitoring.svg?branch=master)](https://travis-ci.org/Auginte/go-monitoring)

Before committing please run `scripts/code-quality.sh`,
so there will be less discussions "between tabs vs spaces"