#!/usr/bin/env bash

if [ ! -e /go/src/github.com/Auginte/go-monitoring/vendor ]; then
    echo "Cannot compile witout dependencies!"
    echo "Install Glide from: https://github.com/Masterminds/glide"
    echo "And run:"
    echo "  glide install"
    echo ""
    echo "[ERROR] Finished without compiling"
    exit 2
fi

if [ ! -e /go/bin/golint ]; then
    echo "Getting linter (will be cached with next run)..."
    go get -u github.com/golang/lint/golint

    echo "Checking code"
fi

go fmt $(go list github.com/Auginte/go-monitoring/... | grep -v /vendor/)
golint github.com/Auginte/go-monitoring/... | grep -v /vendor/
