#!/usr/bin/env bash

DIR="$( cd "$( dirname "$(dirname "${BASH_SOURCE[0]}" )" )" && pwd )"

docker run -it -v "$DIR":/go/src/github.com/Auginte/go-monitoring/ \
               -v "$DIR/cache":/go/bin \
               -v "$DIR/scripts/raw-code-quality.sh":/custom/raw-code-quality.sh \
               golang:1.7.0 /custom/raw-code-quality.sh