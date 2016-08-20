#!/usr/bin/env bash

DIR="$( cd "$( dirname "$(dirname "${BASH_SOURCE[0]}" )" )" && pwd )"

# will build into bin
docker run -it -u `id -u $USER` -v "$DIR":/go/src/github.com/Auginte/go-monitoring/ -v "$DIR/bin":/go/bin golang:1.7.0 go install github.com/Auginte/go-monitoring/...