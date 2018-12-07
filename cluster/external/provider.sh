#!/bin/bash

function _kubectl() {
    kubectl "$@"
}

function prepare_config() {
    cat >hack/config-provider-external.sh <<EOF
docker_tag=devel
docker_prefix=${DOCKER_PREFIX}
manifest_docker_prefix=${DOCKER_PREFIX}
image_pull_policy=${IMAGE_PULL_POLICY:-Always}
EOF
}

# The external cluster is assumed to be up.  Do a simple check
function external_up() {
    prepare_config
    #_kubectl version >/dev/null
    version=$(_kubectl version)
    if [ $? -ne 0 ]; then
        echo -e "\n*** Unable to reach external cluster.  Please check configuration ***"
        echo -e "*** Type \"kubectl config view\" for current settings               ***\n"
        exit 1
    fi
    echo "external provider detected; version info:"
    echo "$version" | sed 's/^/    /'
}

function down() {
    echo "Not supported by this provider"
}

function build() {
    # Build code and manifests
    ${KUBEVIRT_PATH}hack/dockerized "DOCKER_TAG=${DOCKER_TAG} KUBEVIRT_PROVIDER=${KUBEVIRT_PROVIDER} ./hack/build-manifests.sh"
    make push
}
