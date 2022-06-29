[//]: # "SPDX-License-Identifier: Apache-2.0"
[//]: # "Copyright (c) 2020-2022 Intel Corporation"

# Running Test Cases with emcoctl

This folder contains everything you need to run a sample test case with the Temporal Action Controller (TAC). EMCO, temporal, a temporal worker, a temporal client, and the temporal client server need to be running and listening for these tests to work. This document will go over how to set up all the temporal entities, so you can run the test case.

## Bringing Up Temporal for temporal action controller


1. First step is to bring up the temporal server itself. The most common way to bring up a temporal worker is with docker-compose provided by [temporal themselves](https://docs.temporal.io/clusters/quick-install/)

 * Clone the temporal server repo somewhere locally on your machine

 `` git clone https://github.com/temporalio/docker-compose.git ``

 * cd into the directory

  `` cd  docker-compose ``

 * Bring up temporal in detatched mode so it will run in the background

  `` docker-compose -f docker-compose.yml up -Vd ``

2. Once temporal is up and running we need to bring up the temporal client, http server, and temporal worker. In a future release of EMCO the temporal client and worker will be included in this directory, so they can be compiled and run alongside this demo.


## Setup the environment to run test cases

* In the config file in this directory have the following variables set:
    1. ``KUBE_PATH``: points to the absolute path where the kubeconfig for the edge cluster is located.
    2. ``HOST_IP``: IP address of the cluster where EMCO is installed.

* The setup.sh script

    This file creates all of the variables we will need to run all of the yaml test scripts inside of this environment.

    ```
    $ ./setup.sh create
    ```

    Output:
    * ``values.yaml``: specifies useful variables for the creation of EMCO resources


    ```
    $ ./setup.sh cleanup
    ```

    Cleans up all artifacts previously generated.

## Generating the Temporal Environment

## Apply prerequisites to run tests
Apply prerequisites.yaml. This is required for all the tests. This creates controllers, one project, one cluster, a logical cloud. This step is required to be done only once for all usecases:

```
$ emcoctl apply -f 00-prerequisites.yaml -v values.yaml
```

## Initialize logical cloud

```
$ emcoctl apply -f 01-init-lc.yaml -v values.yaml
```

## Define App Deployment Intent Group

```
$ emcoctl apply -f 02-define-app-dig.yaml -v values.yaml
```

## Define workflow intents

```
$ emcoctl apply -f 03-define-workflow-intents.yaml -v values.yaml
```

## Running the test cases

1. The install hooks. The instantiate deployment intent group will run both the pre/post install hooks inside of the database. You can check the logs inside of the workflow client to check them running.


```
$ emcoctl apply -f 04-instantiate-dig.yaml -v values.yaml
```

2. The update hooks. The update deployment intent group will run both the pre/post update hooks.

```
$ emcoctl apply -f 05-update-dig.yaml -v values.yaml
```

3. Finally the delete hooks. Deleting an instantiated deployment intent group will run the pre/post delete hooks inside of the database.

```
$ emcoctl apply -f 06-delete-dig.yaml -v values.yaml
```

## Cleaning up

Run steps 03-00 using delete instead of apply with emcoctl

```
$ emcoctl delete -f 03-define-workflow-intents.yaml -v values.yaml
$ emcoctl delete -f 02-define-app-dig.yaml -v values.yaml
$ emcoctl delete -f 01-init-lc.yaml -v values.yaml
$ emcoctl delete -f 00-prerequisites.yaml -v values.yaml
```

Delete the temporal client and the temporal worker. Navigate back to the helm temporal directory and run these commands.

```
helm uninstall demo1 worker -n demo
helm uninstall demo workflowclient -n demo

kubectl delete ns demo
```

Delete the temporal server itself by navigating to wherever you cloned the docker-compose temporal repository and run:

```
docker-compose down
```


Finally, run the cleanup function from the setup.sh file to cleanup the values.yaml file.

```
$ ./setup.sh cleanup
```