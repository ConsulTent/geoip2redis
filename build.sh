#!/bin/bash
export GOPATH=$HOME/go

if [ "${GOOS}" == "windows" ]; then
go build -v -ldflags="-X main.gitver=$(git describe --always --abbrev=4)" -o bin/geoip2redis.exe *.go
go build -v -ldflags="-X main.gitver=$(git describe --always --long)" -o bin/ip2long.exe tools/ip2long/ip2long.go
go build -v -ldflags="-X main.gitver=$(git describe --always --long)" -o bin/maxmind-ip2location-to-csv.exe tools/maxmind-ip2location-to-csv/maxmind-ip2location-to-csv.go
echo 'bin/geoip2redis.exe bin/ip2long.exe bin/maxmind-ip2location-to-csv.exe'
else
go build -v -ldflags="-X main.gitver=$(git describe --always --abbrev=4)" -o bin/geoip2redis *.go
go build -v -ldflags="-X main.gitver=$(git describe --always --abbrev=4)" -o bin/ip2long tools/ip2long/ip2long.go
go build -v -ldflags="-X main.gitver=$(git describe --always --abbrev=4)" -o bin/maxmind-ip2location-to-csv tools/maxmind-ip2location-to-csv/maxmind-ip2location-to-csv.go
echo 'bin/geoip2redis bin/ip2long bin/maxmind-ip2location-to-csv'
fi
