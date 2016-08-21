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

go install $(go list github.com/Auginte/go-monitoring/... | grep -v /vendor/)
