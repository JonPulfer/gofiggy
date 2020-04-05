#!/bin/sh
set -o errexit

go test -v ./...

docker build -t localhost:5000/gofiggy .
docker tag localhost:5000/gofiggy localhost:5000/gofiggy:latest
docker image prune -f --filter label=stage=builder
docker push localhost:5000/gofiggy
