#!/usr/bin/env bash

#
# Copyright 2019 The Tekton Authors
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
#

set -o errexit
set -o nounset

ORG="github.com/hyperfunction"
PROJECT="hyperfunction"

# Conditionally create a temporary GOPATH for codegen
# and openapigen to execute in. This is only done if
# the current repo directory is not within GOPATH.
function shim_gopath() {
  local REPO_DIR=$(git rev-parse --show-toplevel)
  local TEMP_GOPATH="${REPO_DIR}/.gopath"
  local TEMP_ORG="${TEMP_GOPATH}/src/${ORG}"
  local TEMP_PROJECT="${TEMP_ORG}/${PROJECT}"
  local NEEDS_MOVE=1

  # Checks if GOPATH exists without triggering nounset panic.
  EXISTING_GOPATH=${GOPATH:-}

  # Check if repo is in GOPATH already and return early if so.
  # Unfortunately this doesn't respect a repo that's symlinked into
  # GOPATH and will create a temporary anyway. I couldn't figure out
  # a way to get the absolute path to the symlinked repo root.
  if [ -n "$EXISTING_GOPATH" ] ; then
    case $REPO_DIR/ in
      $EXISTING_GOPATH/*) NEEDS_MOVE=0;;
      *) NEEDS_MOVE=1;;
    esac
  fi

  if [ $NEEDS_MOVE -eq 0 ]; then
    return
  fi

  echo "You appear to be running from outside of GOPATH."
  echo "This script will create a temporary GOPATH at $TEMP_GOPATH for code generation."

  # Ensure that the temporary project symlink doesn't exist before proceeding.
  delete_repo_symlink

  mkdir -p "$TEMP_ORG"
  # This will create a symlink from
  # (repo-root)/.gopath/src/github.com/hyperfunction/hyperfunction
  # to the user's checkout.
  ln -s "$REPO_DIR" "$TEMP_ORG"
  echo "Moving to $TEMP_PROJECT"
  cd "$TEMP_PROJECT"
  export GOPATH="$TEMP_GOPATH"
}

# Helper that wraps deleting the temp repo symlink
# and prints a message about deleting the temp GOPATH.
#
# Why doesn't this func just delete the temp GOPATH outright?
# Because it might be reused across multiple hack scripts and many
# packages seem to be installed into GOPATH with read-only
# permissions, requiring sudo to delete the directory. Rather
# than surprise the dev with a password entry at the end of the
# script's execution we just print a message to let them know.
function shim_gopath_clean() {
  local REPO_DIR=$(git rev-parse --show-toplevel)
  local TEMP_GOPATH="${REPO_DIR}/.gopath"
  if [ -d "$TEMP_GOPATH" ] ; then
    # Put the user back at the root of the project repo
    # after all the symlink shenanigans.
    echo "Moving to $REPO_DIR"
    cd "$REPO_DIR"
    delete_repo_symlink
    echo "When you are finished with codegen you can safely run the following:"
    echo "sudo rm -rf \".gopath\""
   fi
}

# Delete the temp symlink to project repo from the temp GOPATH dir.
function delete_repo_symlink() {
  local REPO_DIR=$(git rev-parse --show-toplevel)
  local TEMP_GOPATH="${REPO_DIR}/.gopath"
  if [ -d "$TEMP_GOPATH" ] ; then
    local REPO_SYMLINK="${TEMP_GOPATH}/src/$ORG/$PROJECT"
    if [ -L $REPO_SYMLINK ] ; then
      echo "Deleting symlink to pipelines repo $REPO_SYMLINK"
      rm -f "${REPO_SYMLINK}"
    fi
  fi
}
