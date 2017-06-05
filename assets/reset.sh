#!/bin/sh
# This script will reset all the database and then create:
# - a new antling with id 1
# - a new image from alpine:latest with id 1
# - a new job with id 1, using image 1 on antling 1
#
# BEWARE: Run it from the root of the project with ./assets/reset.sh
set -xe

go generate cmd/anthivectl.go

go install cmd/anthivectl.go
go install github.com/alienantfarm/anthive

anthivectl reset
anthive &
anthive_pid=$!
rm -f static/images/*
sleep 1 # be sure everything started

curl -X POST -H "Content-Type:application/json" -d '{"tag": "alpine:latest"}' \
	localhost:8888/images
curl -X POST localhost:8888/antlings
curl -X POST -H "Content-Type:application/json" \
	-d '{"image": 1, "command": ["foo", "bar"]}' localhost:8888/jobs

kill $anthive_pid
anthive
