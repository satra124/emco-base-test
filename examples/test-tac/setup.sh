#!/bin/bash

# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2022 Intel Corporation

# Reading the config file
echo "Reading the config file"
CONFIG_HOST_IP=$(cat ./config | grep HOST_IP | cut -d'=' -f2)
CONFIG_KUBE_PATH=$(cat ./config | grep KUBE_PATH | cut -d'=' -f2)
RSYNC_CONTROL_PORT=$(cat ./config | grep RSYNC_CONTROL_PORT | cut -d'=' -f2)
TAC_CONTROL_PORT=$(cat ./config | grep TAC_CONTROL_PORT | cut -d'=' -f2)
WF_CLIENT_NAME=$(cat ./config | grep WF_CLIENT_NAME | cut -d'=' -f2)
WF_PORT=$(cat ./config | grep WF_PORT | cut -d'=' -f2)

# Reading the environment first - it will take prority
HOST_IP=${HOST_IP:-$CONFIG_HOST_IP}
HOST_IP=${HOST_IP:-"oops"}
KUBE_PATH=${KUBE_PATH:-$CONFIG_KUBE_PATH}
KUBE_PATH=${KUBE_PATH:-"oops"}

# Create the values.yaml file to be able to run the TAC Test Case
function create {
echo "Generating values.yaml"
cat << NET >> values.yaml
HostIP: $CONFIG_HOST_IP 
KubeConfig: $CONFIG_KUBE_PATH
RsyncPort: $RSYNC_CONTROL_PORT
TacPort: $TAC_CONTROL_PORT
WfClientName: $WF_CLIENT_NAME
WfClientPort: $WF_PORT
ProjectName: proj1
LogicalCloud: lc1
AdminCloud: default
ClusterProvider: provider1
Cluster1: cluster
ClusterLabel: edge-cluster
Cluster1Ref: cluster1-ref
CompositeApp: capp
DeploymentIntent: dig
GenericPlacementIntent: genericPlacementIntent
TacIntent: generic-intent
NET
}

function cleanup {
    echo "Removing values.yaml"
    rm -f values.yaml
}


case "$1" in 
    "create" )
    if [ "${HOST_IP}" == "oops" ]  || [ "${KUBE_PATH}" == "oops" ] ; then
        echo -e "ERROR - HOST_IP & KUBE_PATH need to be defined as environment variables or in the config file"
    else
        create
    fi
    ;;
    "cleanup" )
    cleanup
    ;;
    *)
esac