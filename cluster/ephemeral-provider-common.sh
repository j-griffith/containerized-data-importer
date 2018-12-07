#!/bin/bash
set -e
function _kubectl() {
 export KUBECONFIG=${CDI_PATH}cluster/$KUBEVIRT_PROVIDER/.kubeconfig
 ${CDI_PATH}cluster/$KUBEVIRT_PROVIDER/.kubectl "$@"
}
