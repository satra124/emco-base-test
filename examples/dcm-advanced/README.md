[//]: # "SPDX-License-Identifier: Apache-2.0"
[//]: # "Copyright (c) 2022 Intel Corporation"

# Advanced DCM (Logical Cloud) operations

The examples herein provide an easy way to test add/removing of clusters to an already-instantiated Logical Cloud.

Here's one possible order of commands to test this:

# Instantiate a logical cloud with one cluster

    emcoctl --config emco-cfg.yaml apply -f prerequisites.yaml -v values.yaml
    emcoctl --config emco-cfg.yaml apply -f 1stcluster.yaml -v values.yaml
    emcoctl --config emco-cfg.yaml apply -f instantiate-lc.yaml -v values.yaml

# Add a cluster reference and update the logical cloud
    emcoctl --config emco-cfg.yaml apply -f 2ndcluster.yaml -v values.yaml
    emcoctl --config emco-cfg.yaml apply -f update-lc.yaml -v values.yaml

# Deploy a composite application to both clusters of the logical cloud
    emcoctl --config emco-cfg.yaml apply -f deployment-intent.yaml -v values.yaml
    emcoctl --config emco-cfg.yaml apply -f instantiate-deployment-intent.yaml -v values.yaml

# Terminate and delete the composite application to both clusters of the logical cloud
    emcoctl --config emco-cfg.yaml delete -f instantiate-deployment-intent.yaml -v values.yaml -w 5
    emcoctl --config emco-cfg.yaml delete -f deployment-intent.yaml -v values.yaml

# Terminate and delete the logical cloud
    emcoctl --config emco-cfg.yaml delete -f instantiate-lc.yaml -v values.yaml -w 5
    emcoctl --config emco-cfg.yaml delete -f 2ndcluster.yaml -v values.yaml
    emcoctl --config emco-cfg.yaml delete -f 1stcluster.yaml -v values.yaml
    emcoctl --config emco-cfg.yaml delete -f prerequisites.yaml -v values.yaml
