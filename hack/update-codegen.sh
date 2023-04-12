#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]:-$0}"; )" &> /dev/null && pwd 2> /dev/null; )";

# Conveniently set GOPATH if unset
if [[ -z "${GOPATH:-}" ]]; then
  export GOPATH="$(go env GOPATH)"
  if [[ -z "${GOPATH}" ]]; then
    echo "WARNING: GOPATH not set and go binary unable to provide it"
  fi
fi

# Useful environment variables
readonly REPO_ROOT_DIR="${REPO_ROOT_DIR:-$(git rev-parse --show-toplevel 2> /dev/null)}"
readonly REPO_NAME="${REPO_NAME:-$(basename ${REPO_ROOT_DIR} 2> /dev/null)}"

# deepcopy and clients generation
bash ${REPO_ROOT_DIR}/hack/generate-groups.sh "deepcopy,client,informer,lister" \
  github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/client github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/apis \
  "pipelines-feedback.keskad.pl:v1alpha1" \
  --go-header-file="${SCRIPT_DIR}/boilerplate.go.txt" \
  --output-base=${SCRIPT_DIR}/../.build/generated

# clean up
rm -rf ${SCRIPT_DIR}/../pkg/client

# copy regenerated
cp -pr ${SCRIPT_DIR}/../.build/generated/github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/apis/* ${SCRIPT_DIR}/../pkg/apis
mv ${SCRIPT_DIR}/../.build/generated/github.com/Kubernetes-Native-CI-CD/pipelines-feedback-core/pkgs/client ${SCRIPT_DIR}/../pkg/client
rm -rf ${SCRIPT_DIR}/../.build/generated
