#!/usr/bin/env bash

docker build -t tinygo .
docker run -it --rm -p 8080:8080 --name tinygo tinygo
