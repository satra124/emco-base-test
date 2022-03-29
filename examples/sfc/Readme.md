#### SPDX-License-Identifier: Apache-2.0
#### Copyright (c) 2021 Intel Corporation

# Running EMCO testcases with emcoctl

This folder contains an example which uses EMCO to deploy
Service Function Chaining (SFC).  The SFC functionality is
provided by edge clusters which use the Nodus CNI (aka OVN4NFV)
https://github.com/akraino-edge-stack/icn-nodus

EMCO provides two action controllers

1. sfc controller - which takes SFC intents 
as part of an EMCO composite application to deploy the
service functions to create an SFC.

2. sfc client controller - takes SFC client intents to
as part of EMCO composite applications which will connect
to either end of the SFC.

## Setup required

1. An Edge cluster which supports Nodus SFC must be prepared.
   - As an example, see instructions found here https://github.com/akraino-edge-stack/icn-nodus/tree/master/demo/calico-nodus-secondary-sfc-setup-II
     - Note that the yaml file used to deploy the Nodus CNI ( https://github.com/akraino-edge-stack/icn-nodus/blob/master/deploy/ovn4nfv-k8s-plugin-sfc-setup-II.yaml ) contains a `ConfigMap` resource called `ovn-controller-network` which is used to define a pool of subnets which are used by Nodus to create virtual networks to stitch together the clients and functions of a network chain.  The element in the configmap used to define the pool of subnets is the `virtual-net-conf.json` element.  On deployment, another configmap is created called `nodus-dynamic-network-pool`, which is used to manage the subnet pool.
   - Ensure the EMCO `monitor` is installed in the edge cluster.
   - Apply the `NetworkAttachmentDefinition` used by EMCO to interact with Nodus in clusters (for applying `Networks`, `ProviderNetworks` and `NetworkChaining` CRs).

```
$ cat ovn-networkobj.yaml
apiVersion: k8s.cni.cncf.io/v1
kind: NetworkAttachmentDefinition
metadata:
  name: ovn-networkobj
  namespace: default
spec:
  config: '{ "cniVersion": "0.3.1", "name": "ovn4nfv-k8s-plugin", "type": "ovn4nfvk8s-cni"
    }'
```

## Description of SFC example using EMCO

This EMCO SFC example is based off of the example in the Nodus repository:
https://github.com/akraino-edge-stack/icn-nodus/tree/master/demo/calico-nodus-secondary-sfc-setup-II#demo-setup
and illustrates how SFC can be deployed and used with EMCO.

The diagram at the above link illustrates quite closely the end result of this EMCO SFC example.

### Key Components of the EMCO SFC Example

This example illustrates a number of EMCO concepts and controllers.  The example also illustrates the full range of
EMCO SFC intents.  For specific use cases, elements can be added or removed. The example is comprised of the following key components:

1. The `prerequisites.yaml` - as with many of the other ECMO examples, sets up some of the basic EMCO resources.  Including
   `controllers`, `clusters`, and, for this example, cluster provider network intents.
1. The `sfc-networks.yaml` is used to instantiate the provider network intents to the actual clusters.  The provider networks
   will be used to illustrate how provider networks can be attached to the ends of a network chain. It is not required to
   use provider networks with network chaining (SFC), so this component and the other intents (`ovnaction`) used for provider networks
   can be removed for specific use cases.
1. The `sfc-ca.yaml` is an EMCO composite application that includes a set of `apps` which comprise the Functions in the SFC.
   There are a few EMCO concepts utilized to set up the SFC with this composite application.
   1. `SFC`, `SFC links`, `SFC client` and `SFC provider network` intents are included to define structure of the SFC.
   1. The `SFC link` intents identify the name of a virtual network that is used to connect elements in the SFC, as well as to
      provide an attachment point for client `Pods` that will be attached to the SFC.  Nodus, via the dynamic subnet pool mechanism,
      mentioned above will automatically create these virtual networks and subnets.  The key link between EMCO and Nodus is the
      name (and label) applied to those networks - which are specified in the `SFC link` intents.
   1. `ovnaction` intents are used to identify which workload resource in the Functions of the composite application will be
      annotated such that they are instantiated with network interfaces on the appropriate virtual networks of the SFC.  These
      `ovnaction` interface intents identify the name of the virtual network and the name of the interface that is to be
      added to the `Pod`.
   1. Also, `ovnaction` intents are used to add provider network interfaces to the appropriate network Functions - i.e. the functions
      at either end of the chain.
1. The `sfc-left-client*.yaml` and `sfc-right-client*.yaml` files provide two clients for the left/head of the SFC and two clients
   for the right/tail of the SFC.  Each of these files are separate composite applications which use the `SFC Client Intent` to
   specify which end of the SFC the identified `Pods` in the application will be connected.
1. The `sfc-tm1.yaml` and `sfc-tm2.yaml` composite applications provide workloads that connect to one of the provider networks
   respectively.  These can be used as in-cluster methods to try out SFC using the provider network attachments to the SFC.
   Of course, the intent of using provider networks with SFC is to be able to attach VMs, servers, etc. that are external to
   the cluster to the SFC.

### SFC Composite Application

An EMCO composite application is created for deploying the SFC.  It is comprised
of the example service functions (SLB, NGFW and SDEWAN) and a set of SFC intents  which define
the following items:

- `SFC Link Intents` define the structure of the chain - i.e. the sequence of the service functions and connecting
  networks
- `SFC Client Selector Intents ` define the namespace and pod label client selector information used to connect client pods to an end of
  the SFC
- `SFC Provider Netowrk Intents` define the provider networks that attach to an end of the SFC

When this SFC composite application is deployed by EMCO, the SFC action controller will
create the appropriate Nodus `networkchaining` CR and deploy it to the edge cluster(s) along
with the service functions.

### SFC Client Composite Applications

Four EMCO composite applications have been created to attach to the ends of the SFC. The composite applications `sfc-left-client.yaml`
and `sfc-right-client.yaml` should be applied before `sfc-left-client-2.yaml` and `sfc-right-client-2.yaml`.  It is not necessary to
apply all four of these client applications.  The purpose of including four is to illustrate that is possible to define more than one
client for either end of the SFC.

Each of these client composite applications have the following:

- a  logical cloud which defines the Kubernetes namespace to be used - each client in this example uses a different logical cloud/namespace (though not required by Nodus SFC)
- the client workload - which in this case - deploys a replicaset of three pods
- an `SFC Client Intent` which identifies the SFC and which end of the chain to connect with

### SFC Provider Network Applications

Two EMCO composite applications have been created to attach to one of the provider networks that are connected to either end of the
SFC.

## Running the SFC example

All scripts and files for the SFC demo are located within this directory and sub-directories.

### Setup the demo

```
./setup.sh cleanup
./setup.sh create
emcoctl --config emco-cfg.yaml apply -f prerequisites.yaml -v values.yaml
emcoctl --config emco-cfg.yaml apply -f sfc-networks.yaml -v values.yaml
emcoctl --config emco-cfg.yaml apply -f sfc-logical-clouds.yaml -v values.yaml
```

These steps clean up any previous demo executions and create the emco-cfg.yaml, values.yaml
and the helm charts used by the demo.

The provider networks are instantiated to the cluster.

The four logical clouds that will
be used by the (up to four) client composite applications are instantiated.  Also, an `admin` logical
cloud is instantiated.  This is used to deploy the SFC Functions and SFC `networkchaining` CR.

### Apply the SFC Composite Application

```
emcoctl --config emco-cfg.yaml apply -f sfc-ca.yaml -v values.yaml
```

The file `sfc-ca.yaml` defines one SFC Provider Network Intent and one SFC Client
Selector Intent for each end of the chain.


### Deploy the SFC Client applications

Deploy the client applications as follows:

```
emcoctl --config emco-cfg.yaml apply -f sfc-left-client.yaml -v values.yaml
emcoctl --config emco-cfg.yaml apply -f sfc-right-client.yaml -v values.yaml
```

If two client applications are wanted on either end, apply one or both of the second clients:

```
emcoctl --config emco-cfg.yaml apply -f sfc-left-client-2.yaml -v values.yaml
emcoctl --config emco-cfg.yaml apply -f sfc-right-client-2.yaml -v values.yaml
```

Note:  these second client YAML files just define an additional deployment intent group using the
same underlying composite applications - which is why they need to be applied after the first set of clients.

### Testing the SFC

After the SFC composite applications have been deployed as described above, the SFC can be tested by
running a traceroute command from a left side client pod to a right side client pod as follows:
In the following example, the three right nginx pods have IP addresses of `172.30.18.[456]`.

```
$ kubectl -n sfc-head exec `kubectl get pods -lsfc=head -n sfc-head  --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}' | head -1` -- traceroute -n -q 1 -I 172.30.18.6
traceroute to 172.30.18.6 (172.30.18.6), 30 hops max, 46 byte packets
 1  172.30.17.3  0.009 ms
 2  172.30.19.3  0.005 ms
 3  172.30.16.3  0.004 ms
 4  172.30.18.6  1.143 ms

```

Note:  the target IP address of clients on the right side may differ depending on which subnet Nodus assigns to the right end virtual network.  Use a command like one of the following
to identify the IP address:

```
kubectl -n sfc-tail exec `kubectl get pods -lsfc=tail -n sfc-tail --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}'` -- ifconfig sn0

If the second set of clients have been deployed:
kubectl -n sfc-tail-two exec `kubectl get pods -lsfc=tail2 -n sfc-tail-two --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}'` -- ifconfig sn0
```

A reverse traceroute - following from the tail of the chain to the head can also be done.  Again, finding the IP addresses used by the Pods in the head application will likely be necessary.
In the following example, the three left nginx pods have IP addresses of `172.30.17.[456]`.

```
$ kubectl -n sfc-tail exec `kubectl get pods -lsfc=tail -n sfc-tail  --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}' | head -1` -- traceroute -n -q 1 -I 172.30.17.4
traceroute to 172.30.17.4 (172.30.17.4), 30 hops max, 46 byte packets
 1  172.30.18.3  0.260 ms
 2  172.30.16.4  0.265 ms
 3  172.30.19.4  0.310 ms
 4  172.30.17.4  0.255 ms
```

### Testing the SFC with Provider Networks
As noted before, the `sfc-tm1.yaml` and `sfc-tm2.yaml` provide in-cluster composite applications that connect the the left and right provider networks.

```
emcoctl --config emco-cfg.yaml apply -f sfc-tm1.yaml -v values.yaml
emcoctl --config emco-cfg.yaml apply -f sfc-tm2.yaml -v values.yaml
```

For now, once these applications are set up, some configuration is required inside each of the pods.

For `sfc-tm1`, exec into the pod and do the following:

```
$ kubectl exec tm1-nginx-7f8477bc8-d8vj8 -it -- bash
bash-5.0# ip route del default
bash-5.0# ip route add default via 172.30.10.3
bash-5.0# exit
```

For `sfc-tm2`, exec into the pod and do the following:

```
$ kubectl exec tm2-nginx-7bbcdd8dc5-gt8kk -it -- bash
bash-5.0# ip route add 172.30.10.0/24 via 172.30.20.3
bash-5.0# ip route add 172.30.16.0/24 via 172.30.20.3
bash-5.0# ip route add 172.30.17.0/24 via 172.30.20.3
bash-5.0# ip route add 172.30.18.0/24 via 172.30.20.3
bash-5.0# ip route add 172.30.19.0/24 via 172.30.20.3
bash-5.0# exit
```

Do a traceroute from the left provider network to a right side client pod
In the following example, the three right nginx pods have IP addresses of `172.30.18.[456]`. The left `sfc-tm1` pod is connected to the left provider network
on the `172.30.10.0` network.  The Function on the left end of the SFC - the `slb` - has IP address `172.30.10.3` on the left provider network.

```
$ kubectl exec `kubectl get pod -l sfc=tm1-nginx --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}'` -it --  traceroute -n -q 1 -I 172.30.18.6
traceroute to 172.30.18.6 (172.30.18.6), 30 hops max, 46 byte packets
 1  172.30.10.3  0.476 ms
 2  172.30.19.3  0.306 ms
 3  172.30.16.3  0.250 ms
 4  172.30.18.6  0.243 ms
```

Do a traceroute from the left provider network to the right provider network.
`172.30.20.120` is the IP address of the `sfc-tm2` pod on the right provider network.

```
$ kubectl exec `kubectl get pod -l sfc=tm1-nginx --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}'` -it --  traceroute -n -q 1 -I 172.30.20.120
traceroute to 172.30.20.120 (172.30.20.120), 30 hops max, 46 byte packets
 1  172.30.10.3  0.570 ms
 2  172.30.19.3  0.415 ms
 3  172.30.16.3  0.286 ms
 4  172.30.20.120  0.313 ms
```

## How to clean up the edge cluster

Delete the emcoctl demo files in reverse order.
The `-w 3` option is useful when deleting - causing `emcoctl` to pause for several seconds after  `terminate` operations - allowing the termination operation to complete
before continuing to delete the resource.

The following shows all the components - skip deleting any of the optional components that may not have been applied (i.e. clients or provider network applications).

```
emcoctl --config emco-cfg.yaml delete -f sfc-tm2.yaml -v values.yaml -w 3
emcoctl --config emco-cfg.yaml delete -f sfc-tm2.yaml -v values.yaml -w 3
emcoctl --config emco-cfg.yaml delete -f sfc-right-client-2.yaml -v values.yaml -w 3
emcoctl --config emco-cfg.yaml delete -f sfc-left-client-2.yaml -v values.yaml -w 3
emcoctl --config emco-cfg.yaml delete -f sfc-right-client.yaml -v values.yaml -w 3
emcoctl --config emco-cfg.yaml delete -f sfc-left-client.yaml -v values.yaml -w 3
emcoctl --config emco-cfg.yaml delete -f sfc-ca.yaml -v values.yaml -w 3
emcoctl --config emco-cfg.yaml delete -f sfc-networks.yaml -v values.yaml -w 3
emcoctl --config emco-cfg.yaml delete -f sfc-logical-clouds.yaml -v values.yaml -w 3
emcoctl --config emco-cfg.yaml delete -f prerequisites.yaml -v values.yaml -w 3
```


## More advanced clean up

Sometimes It may be helpful to clean up and reinstall the Nodus CNI.  The following steps may help accomplish this.

Note: in the examples below, the Nodus repository has been pulled into `/home/vagrant/git`.

```
# after cleaning up the SFC usecase, as described above

# delete the dynamically created virtual networks
kubectl get networks  # this will show what exists

# for this SFC example, the following virtual network get created
kubectl delete network virtual-net1
kubectl delete network virtual-net2
kubectl delete network dynamic-net1
kubectl delete network dynamic-net2

# delete the Nodus CNI
kubectl delete -f /home/vagrant/git/icn-nodus/deploy/ovn4nfv-k8s-plugin-sfc-setup-II.yaml
# wait for above to completely terminate

# delete the dynamic pool subnet configmap
kubectl -n kube-system delete cm nodus-dynamic-network-pool

# delete ovn
kubectl delete -f /home/vagrant/git/icn-nodus/deploy/ovn-daemonset.yaml
# wait for above to completely terminate

# delete multus
kubectl delete -f /home/vagrant/git/icn-nodus/deploy/multus-daemonset.yaml
# wait for above to completely terminate

# after above two are completely terminated, re-apply them
kubectl apply -f /home/vagrant/git/icn-nodus/deploy/multus-daemonset.yaml
# wait for multus to come up

# apply the network attachment definition used by EMCO
kubectl apply -f /home/vagrant/git/icn-nodus/example/emco-net-attach-def-cr.yaml

# apply ovn
kubectl apply -f /home/vagrant/git/icn-nodus/deploy/ovn-daemonset.yaml
# wait for ovn to come up

# apply nodus CNI
kubectl apply -f /home/vagrant/git/icn-nodus/deploy/ovn4nfv-k8s-plugin-sfc-setup-II.yaml
# wait for ovn4nfv to come up
```

## Debugging

The `collect-info.sh` script can be used to collect logs and information about the SFC demo.

> **NOTE** The following manual SFC files were not tested for the `22.03` release

The SFC demo can also be deployed manually (not using EMCO) on the target edge cluster using the instructions
here https://github.com/akraino-edge-stack/icn-nodus/tree/master/demo/calico-nodus-secondary-sfc-setup-II
if needed to verify the edge cluster has been set up correctly.  A copy of the demo SFC resource yaml files
are included in the kud/sfc/manualSFC directory.
