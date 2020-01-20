#!/bin/bash
export GOPATH=$HOME/go
go build -i -v -ldflags="-X main.gitver=$(git describe --always --long --dirty)" -o geoip2redis *.go
go build -i -v -ldflags="-X main.gitver=$(git describe --always --long)" -o tools/ip2long/ip2long tools/ip2long/ip2long.go
go build -i -v -ldflags="-X main.gitver=$(git describe --always --long)" -o tools/maxmind-ip2location/maxmind-ip2location tools/maxmind-ip2location/maxmind-ip2location.go
