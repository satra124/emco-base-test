```text
SPDX-License-Identifier: Apache-2.0
Copyright (c) 2022 Intel Corporation
```
<!-- omit in toc -->
# Running EMCO dependency testcase with Helm Hook support (using emcoctl)
For testing inter app dependency with Helm hook using ngnix, collectd and operatior apps using 2 clusters.

Test Setup has 3 apps and 2 clusters with following setup:
 * collectd (Cluster1 and Cluster2)
 * Operator (Cluster1)
 * Ngnix (Cluster1 and Cluster2). This chart has hooks.

The test setup has following App Dependency:
 * collectd --> Operator --> Ngnix

## Configure
(1) Set the KUBE_PATH1 and KUBE_PATH2 environment variables to cluster kubeconfig file path.

(2) Set the HOST_IP enviromnet variables to the address where emco is running.

(3) Modify examples/test-dependency/setup.sh files to change controller port numbers if they are diffrent than your emco installation.

## Install the application
Install the app using the commands:
```shell
$ cd examples/test-dependency/
$ ./apply.sh
```

## Expected outcome
* Nginx app will come up on both clusters. For the Nginx App there are 2 hooks: Pre-install and Post-install. Pre-install should complete and then the main resources including Nginx pod will comeup. After the Pod is in running state there is a wait time of 10 secs as specified in apps.yaml.
* Operator will next come up on cluster1. 4 Pods come up for this app. After all the Pods are up, there is a wait time of 10 sec as specified in the apps.yaml.
* Next collectd app will come up on cluster1 and cluster2

## Uninstall the application
Uninstall the app using the command:
```shell
$ ./delete.sh
```
