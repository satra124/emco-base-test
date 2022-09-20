#!/bin/bash

# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2022 Intel Corporation

set -o errexit
set -o nounset
set -o pipefail


source ../multi-cluster/_common.sh

test_folder=../../tests/
demo_folder=../../demo/
deployment_folder=../../../deployments/

HOST_IP=${HOST_IP:-"oops"}
KUBE_PATH1=${KUBE_PATH1:-"oops"}
GIT_USER=${GIT_USER:-"oops"}
GIT_TOKEN=${GIT_TOKEN:-"oops"}
GIT_REPO=${GIT_REPO:-"oops"}
OUTPUT_DIR=output
DELAY=${DELAY:-"oops"}
GIT_BRANCH=${GIT_BRANCH:-"oops"}

function create_common_values {
    local output_dir=$1
    local host_ip=$2

    create_apps $output_dir
    create_config_file $host_ip

        cat << NET > values.yaml
    PackagesPath: $output_dir
    ProjectName: proj-anthos
    ClusterProvider: provider-anthos
    ClusterLabel: edge-cluster
    AdminCloud: default
    CompositeApp: test-composite-app
    CompositeProfile: test-composite-profile
    GenericPlacementIntent: test-placement-intent
    DeploymentIntent: test-deployment-intent
    Intent: intent
    RsyncPort: 30431
    GacPort: 30433
    OvnPort: 30432
    DtcPort: 30448
    NpsPort: 30438
    HostIP: $host_ip
    APP1: http-server
    APP2: collectd
    APP3: operator
    GitObj: GitObject
    GitUser: $GIT_USER
    GitToken: $GIT_TOKEN
    GitRepo: $GIT_REPO
    Branch: $GIT_BRANCH
    GitResObj: GitResObject
    ConfigDeleteDelay: $DELAY


    Clusters:
      - Name: cluster2
    Applist:
      - Name: collectd
        Cluster:
          - cluster2

NET
}


function usage {
    echo "Usage: $0 create|cleanup|genhelm"
}

function cleanup {
    rm -f *.tar.gz
    rm -f values.yaml
    rm -f emco-cfg.yaml
    rm -rf $OUTPUT_DIR
}

case "$1" in
    "create" )
        if [ "${HOST_IP}" == "oops" ] ; then
            echo -e "ERROR - Environment variable HOST_IP must be set"
            exit
        fi
        if [ "${GIT_TOKEN}" == "oops"  ] ; then
            echo -e "ERROR - GIT_TOKEN must be defined"
            exit
        fi
        if [ "${GITB_REPO}" == "oops"  ] ; then
            echo -e "ERROR - GIT_REPO must be defined"
            exit
        fi
        if [ "${GIT_USER}" == "oops" ] ; then
            echo -e "GIT_USER must be defined"
        else
            create_common_values $OUTPUT_DIR $HOST_IP
            echo "Done create!!!"
        fi
        ;;
    "cleanup" )
        cleanup
    ;;
    "genhelm" )
	create_apps $OUTPUT_DIR
        echo "Generated Helm chart packages in output/ directory!"
    ;;
    *)
        usage ;;
esac
