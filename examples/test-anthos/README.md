[//]: # "SPDX-License-Identifier: Apache-2.0"
[//]: # "Copyright (c) 2020-2022 Intel Corporation"

# Running Google Anthos testcase with emcoctl

This folder contains the collectd testcase to be run with EMCO deployed to Google Anthos clusters. This test assumes that a GKE cluster with Anthos Config Management enabled has been created as mentioned in https://gitlab.com/project-emco/core/emco-base/-/tree/main/docs/design/gitops_support.md. See https://cloud.google.com/anthos-config-management/docs for Google Anthos-specific documentation.

## Organize git repository structure to fit EMCO and Anthos' needs
Coming later

## Installing monitor

1. fixup monitor helm chart to remove status from RBS CRD template file
2. ideally should release a monitor 1.0.1 after this
3. follow steps below:

```
cd ~/emco-base/deployments/helm
helm package monitor
tar -xf monitor-1.0.0.tgz
CLUSTER_REF="provider-anthos+cluster2"
GITHUB_OWNER="igordcard"
GITHUB_TOKEN=""
GITHUB_REPO="anthosync"
helm template emco monitor -n emco --set git.token=$GITHUB_TOKEN --set git.repo=$GITHUB_REPO --set git.username=$GITHUB_OWNER --set git.clustername=$CLUSTER_REF --set git.enabled=true > monitor.yaml
cp monitor.yaml ~/anthosync/rootsync/acm-cluster/
cd ~/anthosync/rootsync/acm-cluster/
git add monitor.yaml
git commit -a
git push
```

**Note:** following the steps above for a public GitHub repo will expose the access token to the Internet, thus making the entire deployment vulnerable. Either make sure the repository is private, or install monitor via another method such as a direct `kubectl apply -f [prefix/]monitor.yaml` on the GKE clusters.

## Setup Environment variables

Setup environment variables as mentioned in https://gitlab.com/project-emco/core/emco-base/-/tree/main/docs/design/gitops_support.md.

## Generating files

Creates artifacts needed to run the testcase.

```
$ cd emco-base/examples/test-anthos
```
```
$ ./setup.sh create
```

## Creating the controllers

Apply 00-controllers.yaml, this creates the controllers. This step is required to be done only once.

```
$ emcoctl --config emco-cfg.yaml apply -f 00-controllers.yaml -v values.yaml
```

## Applying the prerequisites

Apply 01-prerequisites.yaml.

```
$ emcoctl --config emco-cfg.yaml apply -f 01-prerequisites.yaml -v values.yaml
```

## Create the deployment intent

Apply 02-deployment-intent.yaml.

```
$ emcoctl --config emco-cfg.yaml apply -f 02-deployment-intent.yaml -v values.yaml
```

## Running the test case

```
$ emcoctl --config emco-cfg.yaml apply -f 03-instantiation.yaml -v values.yaml
```

After this step, we should see the resources created in the Google Anthos cluster.


## Cleanup

1. Delete usecase

    ```
    $ emcoctl --config emco-cfg.yaml delete -f 03-instantiation.yaml -v values.yaml
    ```

2. Cleanup deployment intent.

    ```
    $ emcoctl --config emco-cfg.yaml delete -f 02-deployment-intent.yaml -v values.yaml
    ```

3. Cleanup prerequisites

    ```
    $ emcoctl --config emco-cfg.yaml delete -f 01-prerequisites.yaml -v values.yaml
    ```

4. Cleanup controllers

    ```
    $ emcoctl --config emco-cfg.yaml delete -f 00-controllers.yaml -v values.yaml
    ```

5. Cleanup generated files

    ```
    $ ./setup.sh cleanup
    ```
