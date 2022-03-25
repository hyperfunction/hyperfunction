#!/usr/bin/env bash

# Copyright 2019 The Knative Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

REPO_ROOT_DIR=${REPO_ROOT_DIR:-$(git rev-parse --show-toplevel)}

# Run a go tool, installing it first if necessary.
# Parameters: $1 - tool package/dir for go install.
#             $2 - tool to run.
#             $3..$n - parameters passed to the tool.
function run_go_tool() {
  local tool=$2
  if [[ -z "$(which ${tool})" ]]; then
    go install $1
  fi
  shift 2
  ${tool} "$@"
}

# Update licenses.
# Parameters: $1 - output file, relative to repo root dir.
#             $2...$n - directories and files to inspect.
function update_licenses() {
  cd ${REPO_ROOT_DIR} || return 1
  local dst=$1
  shift

  run_go_tool github.com/google/go-licenses go-licenses \
    save ./... --save_path=${dst} --force
  # Hack to make sure directories retain write permissions after save. This
  # can happen if the directory being copied is a Go module.
  # See https://github.com/google/go-licenses/issues/11
  chmod +w $(find ${dst} -type d)
}

echo "Prune modules"
go mod tidy -compat=1.17
go mod vendor

echo "Updating licenses"
update_licenses third_party/VENDOR-LICENSE
