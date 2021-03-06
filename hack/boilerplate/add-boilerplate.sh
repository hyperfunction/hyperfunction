#!/usr/bin/env bash

# Copyright 2022 The hyperfunction Authors
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

USAGE=$(cat <<EOF
Add boilerplate.<ext>.txt to all .<ext> files missing it in a directory.

Usage: (from repository root)
       ./hack/boilerplate/add-boilerplate.sh <ext> <DIR>

Example: (from repository root)
         ./hack/boilerplate/add-boilerplate.sh go cmd
EOF
)

set -e

if [[ -z $1 || -z $2 ]]; then
  echo "${USAGE}"
  exit 1
fi

function grep() {
  local tool=grep
  # Fix compat with mac
  # Install homebrew grep package with `brew install grep`
  [[ -n "$(which ggrep)" ]] && tool=ggrep
  $tool "$@"
}

grep -r -L -P "Copyright \d+ The \w+ Authors" $2  \
  | grep -P "\.$1\$" \
  | xargs -I {} sh -c \
  "cat hack/boilerplate/boilerplate.$1.txt > /tmp/boilerplate && echo '\n' >> /tmp/boilerplate && cat {} >> /tmp/boilerplate && cat /tmp/boilerplate > {}"
