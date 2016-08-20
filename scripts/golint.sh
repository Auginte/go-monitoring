#!/usr/bin/env bash

if [ ! -e /go/bin/golint ]; then
    echo "Getting linter (will be cached with next run)..."
    go get -u github.com/golang/lint/golint

    echo "Checking code"
fi

golint github.com/Auginte/go-monitoring/...
