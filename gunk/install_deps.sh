#!/usr/bin/env bash

set -euo pipefail

source ./devenv.sh

GOBIN=$gunk_dir/bin go install github.com/gunk/gunk
