#!/usr/bin/env bash

DIR="$( cd "$( dirname "$(dirname "${BASH_SOURCE[0]}" )" )" && pwd )"

# will build into bin
docker run -it -v "$DIR":/go/src/github.com/Auginte/go-monitoring/ \
               -v "$DIR/bin":/go/bin \
               -v "$DIR/scripts/raw-build.sh":/custom/raw-build.sh \
               golang:1.7.0 /custom/raw-build.sh