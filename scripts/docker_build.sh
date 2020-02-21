#!/bin/bash

set -euo pipefail

poetry export -f requirements.txt -o requirements.txt
docker build -t armakuni/circleci-workflow-dashboard .
