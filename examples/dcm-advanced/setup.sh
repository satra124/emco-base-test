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
KUBE_PATH2=${KUBE_PATH2:-"oops"}
OUTPUT_DIR=output

function create_common_values {
    local output_dir=$1
    local host_ip=$2

    create_apps $output_dir
    create_config_file $host_ip

        cat << NET > values.yaml
    PackagesPath: $output_dir
    ProjectName: proj1
    ClusterProvider: provider1
    Cluster1: cluster1
    Cluster2: cluster2
    Cluster1Ref: cluster1-ref
    Cluster2Ref: cluster2-ref
    LogicalCloud: lc1
    AdminCloud: admin-lc
    PrivilegedCloud: privileged-lc
    StandardCloud: standard-lc
    PrivilegedNamespace: privileged-lc-ns
    StandardNamespace: standard-lc-ns
    ClusterQuota: quota1
    StandardPermission: standard-permission
    PrivilegedPermission1: privileged-permission-primary
    PrivilegedPermission2: privileged-permission-kube-system
    PrivilegedPermission3: privileged-permission-cluster
    CompositeApp: test-collectd
    HelmApp1: output/collectd.tar.gz
    ProfileApp1: output/collectd_profile.tar.gz
    CompositeProfile: collectd-composite-profile
    GenericPlacementIntent: collectd-placement-intent
    DeploymentIntent: collectd-deployment-intent-group
    Intent: intent
    ConfigmapFile: info.json
    RsyncPort: 30431
    GacPort: 30433
    OvnPort: 30432
    DtcPort: 30448
    NpsPort: 30438
    HostIP: $host_ip

    Clusters:
      - Name: cluster1
        KubeConfig: $KUBE_PATH1
      - Name: cluster2
        KubeConfig: $KUBE_PATH2
    Applist:
      - Name: collectd
        Cluster:
          - cluster1
          - cluster2
NET
}


function usage {
    echo "Usage: $0 create|cleanup"
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
            echo -e "ERROR - Environment varaible HOST_IP must be set"
            exit
        fi
        if [ "${KUBE_PATH1}" == "oops"  ] ; then
            echo -e "ERROR - KUBE_PATH1 must be defined"
            exit
        fi
        if [ "${KUBE_PATH2}" == "oops" ] ; then
            echo -e "KUBE_PATH2 must be defined"
        else
            create_common_values $OUTPUT_DIR $HOST_IP
            echo "Done create!!!"
        fi
        ;;
    "cleanup" )
        cleanup
    ;;
    *)
        usage ;;
esac
