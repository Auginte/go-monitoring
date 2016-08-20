Go monitoring tools
===================

*Augitne.com* specific tools for dealing with metrics and logs transformation for `ElasticSearch`/`Kibana`.

Using `Golang` because `Logstash` takes too much resources and `golang` can download/parse/publish logs in parallel.

Building with docker
--------------------

 * Install [docker](https://docs.docker.com/engine/installation/)
 * Run `scripts/build.sh`
 * Check results in `bin` folder
 
Developing with local go
------------------------

 * [Install go 1.7](https://golang.org/doc/install)
 * clone this project into `$GOPATH/src/github.com/Auginte/go-monitoring/`
 * [Run from IDE](https://plugins.jetbrains.com/plugin/5047)

Tests?
------

Currently everything is tested manually: on Ubuntu and Amazon AMI.

Before committing please run `scripts/code-quality.sh`,
so there will be less discussions "between tabs vs spaces"