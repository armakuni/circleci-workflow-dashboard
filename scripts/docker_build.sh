#!/bin/bash

set -euo pipefail

go build
docker build -t armakuni/circleci-workflow-dashboard .
