#!/bin/bash
export GOPATH=$HOME/go
vgo build -i -v -ldflags="-X main.gitver=$(git describe --always --long --dirty)" -o geoip2redis *.go
