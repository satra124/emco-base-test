#!/bin/bash
# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2020-2022 Intel Corporation


test_case_name=$1
action=$2


function create {
    ./setup.sh create
}

function cleanup {
    ./setup.sh cleanup
}

function get_variables {
    project=$(cat values.yaml | grep ProjectName: | sed -z 's/.*ProjectName: //')
    logical_cloud_name=$(cat values.yaml | grep LogicalCloud: | sed -z 's/.*LogicalCloud: //')

    case $test_case_name in
        "test-prometheus-collectd")
            composite_app="prometheus-collectd-composite-app"
            deployment_intent_group_name="collection-deployment-intent-group"
            ;;
        "test-dtc")
            composite_app="dtc-composite-app"
            deployment_intent_group_name="collection-deployment-intent-group"
            ;;
        "test-vfw")
            composite_app="compositevfw"
            deployment_intent_group_name="vfw_deployment_intent_group"
            ;;
        "monitor")
            composite_app="monitor-composite-app"
            deployment_intent_group_name="collection-deployment-intent-group"
            ;;
        *)
            echo "Invalid testcase file"
            exit
            ;;
    esac
}

function apply_prerequisites {
    emcoctl --config emco-cfg.yaml apply -f prerequisites.yaml -v values.yaml -s
}

function apply_logical_cloud {
    emcoctl --config emco-cfg.yaml apply -f instantiate-lc.yaml -v values.yaml -s
}

function apply_deployment {
    emcoctl --config emco-cfg.yaml apply -f "${test_case_name}-deployment.yaml" -v values.yaml -s
}

function apply_instantiate_testcase {
    emcoctl --config emco-cfg.yaml apply -f "${test_case_name}-instantiate.yaml" -v values.yaml -s
}

function delete_prerequisites {
    emcoctl --config emco-cfg.yaml delete -f prerequisites.yaml -v values.yaml -s
}

function delete_logical_cloud {
    emcoctl --config emco-cfg.yaml delete -f instantiate-lc.yaml -v values.yaml -s
}

function delete_deployment {
    emcoctl --config emco-cfg.yaml delete -f "${test_case_name}-deployment.yaml" -v values.yaml -s
}

function delete_instantiate_testcase {
    emcoctl --config emco-cfg.yaml delete -f "${test_case_name}-instantiate.yaml" -v values.yaml -s
}

# Function to obtain logical cloud status and wait for it to get instantiated
function get_logical_cloud_apply_status {
    echo "Logical Cloud creation in progress... Please Wait"
    for try in {0..300}; do
        sleep 1
        lc_status="$(emcoctl --config emco-cfg.yaml get projects/$project/logical-clouds/$logical_cloud_name/status |  sed -z 's/.*Response://'| jq -r .deployedStatus)"
        case $lc_status in
            "InstantiateFailed")
                echo "Instantiation of logical cloud failed."
                exit
                ;;

            "Instantiating")
                echo -n "."
                continue
                ;;

            "Instantiated")
                echo "Logical cloud creation succeeded!"
                break
                ;;

            *)
                echo "Invalid apply status State."
                exit
                ;;
        esac
    done

    if (($try == 300)) ; then
        echo "Timeout of 300s exceeded!"
    fi

}

# Function to obtain deployment status and wait for it to get instantiated
function get_deployment_intent_group_apply_status {
    echo "Deployment in progress... Please Wait"
    for try in {0..300}; do
        sleep 1
        deployment_status="$(emcoctl --config emco-cfg.yaml get projects/$project/composite-apps/$composite_app/v1/deployment-intent-groups/$deployment_intent_group_name/status?status=deployed |  sed -z 's/.*Response://' | jq -r .deployedStatus)"

        case $deployment_status in
            "InstantiateFailed")
                echo "Instantiation of deployment intent group failed."
                exit
                ;;

            "Instantiating")
                echo -n "."
                continue
                ;;

            "Instantiated")
                echo "Deployment succeeded!"
                break
                ;;

            *)
                echo "Invalid apply status State."
                exit
                ;;
        esac
    done

    if (($try == 300)) ; then
        echo "Timeout of 300s exceeded!"
    fi
}

# Function to obtain logical cloud status and wait for it to get terminated
function get_logical_cloud_delete_status {
    echo "Logical Cloud deletion in progress... Please Wait"
    for try in {0..300}; do
        sleep 1
        lc_status="$(emcoctl --config emco-cfg.yaml get projects/$project/logical-clouds/$logical_cloud_name/status |  sed -z 's/.*Response://' | jq -r .deployedStatus)"

        case $lc_status in
            "TerminateFailed")
                echo "Termination of logical cloud failed."
                exit
                ;;

            "Terminating")
                echo -n "."
                continue
                ;;

            "Terminated")
                echo "Logical Cloud termination succeeded!"
                break
                ;;

            *)
                echo "Invalid delete status State."
                exit
                ;;
        esac
    done

    if (($try == 300)) ; then
        echo "Timeout of 300s exceeded!"
    fi

}

# Function to obtain deployment status and wait for it to get terminated
function get_deployment_intent_group_delete_status {
    echo "Deployment deletion in progress... Please Wait"
    for try in {0..300}; do
        sleep 1
        deployment_status="$(emcoctl --config emco-cfg.yaml get projects/$project/composite-apps/$composite_app/v1/deployment-intent-groups/$deployment_intent_group_name/status?status=deployed |  sed -z 's/.*Response://' | jq -r .deployedStatus)"

        case $deployment_status in
            "TerminateFailed")
                echo "Termination of deployment intent group failed."
                exit
                ;;

            "Terminating")
                echo -n "."
                continue
                ;;

            "Terminated")
                echo "Deployment termination succeeded!"
                break
                ;;

            *)
                echo "Invalid delete status State."
                exit
                ;;
        esac
    done

    if (($try == 300)) ; then
        echo "Timeout of 300s exceeded!"
    fi
}

function usage {

    echo "Usage: $0 <test case file> apply|delete"
    echo "Example: $0 test-prometheus-collectd.yaml apply"

}

function apply {
    cleanup
    create
    get_variables
    apply_prerequisites
    apply_logical_cloud
    get_logical_cloud_apply_status # wait till Instantiated status is obtained
    apply_deployment
    apply_instantiate_testcase
    get_deployment_intent_group_apply_status # wait till Instantiated status is obtained
}

function delete {
    get_variables
    delete_instantiate_testcase
    get_deployment_intent_group_delete_status # wait till Terminated status is obtained
    delete_deployment
    delete_logical_cloud
    get_logical_cloud_delete_status # wait till Terminated status is obtained
    delete_prerequisites
    cleanup
}

if [ "$#" -lt 2 ] ; then
    usage
    exit
fi

case "$2" in
    "apply" ) apply ;;
    "delete" ) delete ;;
    *) usage ;;
esac
