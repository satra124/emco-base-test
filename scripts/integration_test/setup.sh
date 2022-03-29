#!/bin/bash


operator_latest_folder=../../examples/helm_charts/operators-latest
monitor_folder=../../examples/helm_charts/monitor

function create {
    mkdir -p tgz_files
    tar -czf tgz_files/operator_profile.tar.gz -C $operator_latest_folder/profile .
    tar -czf tgz_files/monitor.tar.gz -C $monitor_folder/helm monitor
    rm -f configuration/.env
    touch configuration/.env
    cat << NET >> configuration/.env

HOST=''
ORCHESTRATOR_PORT='30415'
CLM_PORT='30461'
NCM_PORT='30431'
OVNACTION_PORT='30471'
DCM_PORT='30477'
GAC_PORT='30491'
DTC_PORT='30481'
HPA_PORT='30451'
KUBECONFIG_PATH='/home/vagrant/.kube'

NET
}

function cleanup {
    rm -rf tgz_files
}


case $1 in
    create) create ;;
    cleanup) cleanup ;;
    *) echo "Unknown command: $1"
       exit ;;
esac
