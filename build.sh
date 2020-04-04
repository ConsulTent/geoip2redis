#!/bin/bash
export GOPATH=$HOME/go

if [ "${GOOS}" == "windows" ]; then
go build -i -v -ldflags="-X main.gitver=$(git describe --always --long --dirty)" -o geoip2redis.exe *.go
go build -i -v -ldflags="-X main.gitver=$(git describe --always --long)" -o tools/ip2long/ip2long.exe tools/ip2long/ip2long.go
go build -i -v -ldflags="-X main.gitver=$(git describe --always --long)" -o tools/maxmind-ip2location/maxmind-ip2location.exe tools/maxmind-ip2location/maxmind-ip2location.go
echo "geoip2redis.exe tools/ip2long/ip2long.exe tools/maxmind-ip2location/maxmind-ip2location.exe"
else
go build -i -v -ldflags="-X main.gitver=$(git describe --always --long --dirty)" -o geoip2redis *.go
go build -i -v -ldflags="-X main.gitver=$(git describe --always --long)" -o tools/ip2long/ip2long tools/ip2long/ip2long.go
go build -i -v -ldflags="-X main.gitver=$(git describe --always --long)" -o tools/maxmind-ip2location/maxmind-ip2location tools/maxmind-ip2location/maxmind-ip2location.go
echo "geoip2redis tools/ip2long/ip2long tools/maxmind-ip2location/maxmind-ip2location"
fi
