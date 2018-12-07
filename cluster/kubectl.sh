#!/bin/bash

set -e

source $(dirname "$0")/../hack/build/common.sh
source $(dirname "$0")/../hack/build/config.sh

source ${CDI_DIR}/cluster/$KUBEVIRT_PROVIDER/provider.sh
source ${CDI_DIR}/hack/build/config.sh

if [ "$1" == "console" ] || [ "$1" == "vnc" ]; then
    ${CDI_DIR}/_out/cmd/virtctl/virtctl --kubeconfig=${kubeconfig} "$@"
else
    _kubectl "$@"
fi
