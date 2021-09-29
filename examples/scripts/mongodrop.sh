#!/bin/bash
# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2020 Intel Corporation

# Simple script to view appcontext
# with no argumnet, it will list all keys under /context/
# with 1 argument, it will show the value of the key provided
# note: assumes emoco services are running in namespace emco

#### if etcd authentication has been enabled, then modify the command line with the user credentials as shown below:
#    kubectl -n emco exec `kubectl get pods -lapp=etcd -n emco --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}'` -it -- etcdctl --user <userid>:<password> get /context/ --prefix=true --keys-only=true
###

#if [ "$#" -ne 1 ] ; then
#    kubectl -n emco exec `kubectl get pods -lapp=emco-etcd -n emco --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}'` -it -- sh -c "ETCDCTL_API=3 etcdctl get /context/ --prefix=true --keys-only=true"
#else
#if [ "$1" == "del" ] ; then
#    kubectl -n emco exec `kubectl get pods -lapp=emco-etcd -n emco --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}'` -it -- sh -c "ETCDCTL_API=3 etcdctl del /context/ --prefix=true"
#else
#    kubectl -n emco exec `kubectl get pods -lapp=emco-etcd -n emco --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}'` -it -- sh -c "ETCDCTL_API=3 etcdctl get $1 --prefix=true"
#fi

drop="$(cat<<EOF
use emco
db.getCollection('resources').drop()
EOF
)"

kubectl -n emco exec `kubectl get pods -lapp=emco-mongo -n emco --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}'` -it -- sh -c "echo \"$drop\" | mongo"
