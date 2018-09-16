#!/bin/bash

# The only argument this script should ever be called with is '--verify-only'

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[@]}")/..
CODEGEN_PKG="$(ls -d -1 $GOPATH/src/k8s.io/code-generator)"

echo $CODEGEN_PKG

echo "Calling ${CODEGEN_PKG}/generate-groups.sh"
"${CODEGEN_PKG}"/generate-groups.sh all \
  github.com/sbezverk/nsm2nsm/pkg/client github.com/sbezverk/nsm2nsm/pkg/apis \
  sbezverk.io:v1

echo "Generating other deepcopy funcs"
"${GOPATH}"/bin/deepcopy-gen \
  --input-dirs ./pkg/apis/sbezverk.io/v1 \
  --bounding-dirs ./pkg/apis/sbezverk.io/v1 \
  -O zz_generated.deepcopy
