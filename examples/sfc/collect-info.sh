#!/bin/bash
# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2021 Intel Corporation

if [ "$#" -lt 1 ] ; then
  echo "enter a destination directory"      
  exit  
fi

mkdir $1
if [ "$?" -ne 0 ] ; then
	echo "directory already exists"
	exit
fi


kubectl get pod `kubectl get pods -lapp=ngfw --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}'` -o yaml > $1/ngfw-live.yaml
kubectl exec `kubectl get pods -lapp=ngfw --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}'` -- route -n > $1/ngfw.route
kubectl exec `kubectl get pods -lapp=ngfw --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}'` -- ifconfig > $1/ngfw.ifconfig
kubectl get pod `kubectl get pods -lapp=slb --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}'` -o yaml > $1/slb-live.yaml
kubectl exec `kubectl get pods -lapp=slb --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}'` -- route -n > $1/slb.route
kubectl exec `kubectl get pods -lapp=slb --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}'` -- ifconfig > $1/slb.ifconfig
kubectl get pod `kubectl get pods -lapp=sdewan --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}'` -o yaml > $1/sdewan-live.yaml
kubectl exec `kubectl get pods -lapp=sdewan --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}'` -- route -n > $1/sdewan.route
kubectl exec `kubectl get pods -lapp=sdewan --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}'` -- ifconfig > $1/sdewan.ifconfig
kubectl get network virtual-net1 -o yaml > $1/virtual-net1-live.yaml
kubectl get network virtual-net2 -o yaml > $1/virtual-net2-live.yaml
kubectl get network dync-net2 -o yaml > $1/dync-net2-live.yaml
kubectl get network dync-net1 -o yaml > $1/dync-net1-live.yaml
kubectl -n sfc-head get pod/`kubectl get pods -lsfc=head -n sfc-head  --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}' | head -1` -o yaml > $1/nginx-left-live.yaml

kubectl -n sfc-head exec `kubectl get pods -lsfc=head -n sfc-head  --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}' | head -1` -- route -n > $1/nginx-left-1.route
kubectl -n sfc-head exec `kubectl get pods -lsfc=head -n sfc-head  --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}' | head -2 | tail -1` -- route -n > $1/nginx-left-2.route
kubectl -n sfc-head exec `kubectl get pods -lsfc=head -n sfc-head  --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}' | head -3 | tail -1` -- route -n > $1/nginx-left-3.route

kubectl -n sfc-head exec `kubectl get pods -lsfc=head -n sfc-head  --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}' | head -1` -- ifconfig > $1/nginx-left-1.ifconfig
kubectl -n sfc-head exec `kubectl get pods -lsfc=head -n sfc-head  --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}' | head -2 | tail -1` -- ifconfig > $1/nginx-left-2.ifconfig
kubectl -n sfc-head exec `kubectl get pods -lsfc=head -n sfc-head  --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}' | head -3 | tail -1` -- ifconfig > $1/nginx-left-3.ifconfig

kubectl -n sfc-tail get pod/`kubectl get pods -lsfc=tail -n sfc-tail  --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}' | head -1` -o yaml > $1/nginx-right-live.yaml
kubectl -n sfc-tail exec `kubectl get pods -lsfc=tail -n sfc-tail  --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}' | head -1` -- route -n > $1/nginx-right-1.route
kubectl -n sfc-tail exec `kubectl get pods -lsfc=tail -n sfc-tail  --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}' | head -2 | tail -1` -- route -n > $1/nginx-right-2.route
kubectl -n sfc-tail exec `kubectl get pods -lsfc=tail -n sfc-tail  --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}' | head -3 | tail -1` -- route -n > $1/nginx-right-3.route
kubectl -n sfc-tail exec `kubectl get pods -lsfc=tail -n sfc-tail  --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}' | head -1` -- ifconfig > $1/nginx-right-1.ifconfig
kubectl -n sfc-tail exec `kubectl get pods -lsfc=tail -n sfc-tail  --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}' | head -2 | tail -1` -- ifconfig > $1/nginx-right-2.ifconfig
kubectl -n sfc-tail exec `kubectl get pods -lsfc=tail -n sfc-tail  --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}' | head -3 | tail -1` -- ifconfig > $1/nginx-right-3.ifconfig
kubectl get ns sfc-head -o yaml > $1/sfc-head-live.yaml
kubectl get ns sfc-tail -o yaml > $1/sfc-tail-live.yaml
kubectl -n kube-system logs --tail 5000 -l name=nfn-operator  > $1/nfn-operator.log
kubectl -n kube-system logs --tail 5000 -l app=nfn-agent  > $1/nfn-agent.log
kubectl get providernetwork left-pnetwork -o yaml > $1/left-pnetwork-live.yaml
kubectl get providernetwork right-pnetwork -o yaml > $1/right-pnetwork-live.yaml
kubectl get ns sfc-head -o yaml > $1/sfc-head-live.yaml
kubectl get ns sfc-tail -o yaml > $1/sfc-tail-live.yaml
kubectl -n sfc-head exec `kubectl get pods -lsfc=head -n sfc-head  --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}' | head -1` -- traceroute -n -q 1 -I 172.30.22.4 > $1/traceroute.out
