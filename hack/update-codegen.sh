#!/usr/bin/env bash

# Copyright 2022 The hyperfunction Authors
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

source "$(git rev-parse --show-toplevel)/hack/setup-temporary-gopath.sh"
shim_gopath
trap shim_gopath_clean EXIT

source "$(git rev-parse --show-toplevel)/vendor/knative.dev/hack/codegen-library.sh"

# If we run with -mod=vendor here, then generate-groups.sh looks for vendor files in the wrong place.
export GOFLAGS=-mod=

boilerplate="${REPO_ROOT_DIR}/hack/boilerplate/boilerplate.go.txt"

echo "=== Update Codegen for ${MODULE_NAME}"

group "Kubernetes Codegen"

# generate the code with:
# --output-base    because this script should also be able to run inside the vendor dir of
#                  k8s.io/kubernetes. The output-base is needed for the generators to output into the vendor dir
#                  instead of the $GOPATH directly. For normal projects this can be dropped.
${REPO_ROOT_DIR}/hack/generate-groups.sh "deepcopy,client,informer,lister" \
  github.com/hyperfunction/hyperfunction/pkg/client github.com/hyperfunction/hyperfunction/pkg/apis \
  "core:v1alpha1" \
  --go-header-file "${boilerplate}"

group "Knative Codegen"

# Knative Injection
${REPO_ROOT_DIR}/hack/generate-knative.sh "injection" \
  github.com/hyperfunction/hyperfunction/pkg/client github.com/hyperfunction/hyperfunction/pkg/apis \
  "core:v1alpha1" \
  --go-header-file "${boilerplate}"
