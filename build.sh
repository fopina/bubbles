#!/bin/bash

set -e

cd $(dirname $0)

docker pull --platform=arm golang:1.12-stretch
docker build -t fopina/misc:bubbles-arm-builder  .
docker push fopina/misc:bubbles-arm-builder