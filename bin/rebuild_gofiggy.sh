#!/bin/sh
set -o errexit

docker build -t localhost:5000/gofiggy .
docker image prune -f --filter label=stage=builder
docker push localhost:5000/gofiggy
