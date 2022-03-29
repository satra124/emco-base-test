#!/bin/bash

operator_latest_folder=../../examples/helm_charts/operators-latest
monitor_folder=../../deployments/helm

function create {
    mkdir -p tgz_files
    tar -czf tgz_files/operator_profile.tar.gz -C $operator_latest_folder/profile .
    tar -czf tgz_files/monitor.tar.gz -C $monitor_folder monitor
}

function cleanup {
    rm -rf tgz_files
}


case $1 in
    create) create ;;
    cleanup) cleanup ;;
    *) echo "Usage: ./setup.sh create|cleanup"
       exit ;;
esac
