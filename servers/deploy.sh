#!/usr/bin/env bash

./build.sh
echo "Deploying server to Dockerhub"
docker build -t ssharif6/simplify .
docker push ssharif6/simplify

echo "Running Docker instance on Cloud Desktop"
ssh root@simplify.api.shaheensharifian.me 'bash -s' < provision.sh
go clean
