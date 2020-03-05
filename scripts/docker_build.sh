#!/bin/bash

set -euo pipefail

env GOOS=linux GARCH=amd64 CGO_ENABLED=0 go build
docker build -t armakuni/circleci-workflow-dashboard .
