#!/usr/bin/env bash

DIR="$( cd "$( dirname "$(dirname "${BASH_SOURCE[0]}" )" )" && pwd )"

docker run -it -u `id -u $USER` -v "$DIR":/go/src/github.com/Auginte/go-monitoring/ golang:1.7.0 go fmt github.com/Auginte/go-monitoring/...
docker run -v "$DIR":/go/src/github.com/Auginte/go-monitoring/ -v "$DIR/cache":/go/bin -v "$DIR/scripts/golint.sh":"/custom/golint.sh" golang:1.7.0 "/custom/golint.sh"
