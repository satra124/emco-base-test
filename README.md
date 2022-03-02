```
SPDX-License-Identifier: Apache-2.0
Copyright (c) 2019-2022 Intel Corporation
```

# EMCO

## Overview

The Edge Multi-Cluster Orchestrator (EMCO) is a software framework for
intent-based deployment of cloud-native applications to a set of Kubernetes
clusters, spanning enterprise data centers, multiple cloud service providers
and numerous edge locations. It is architected to be flexible, modular and
highly scalable. It is aimed at various verticals, including telecommunication
service providers.

Refer to [EMCO documentation](docs/design/emco-design.md) for details on EMCO architecture.

## Build and Deploy EMCO

### Set up the environment

Depending on your scenario, you may have to set different environment variables.
If behind an HTTP proxy for the Internet, make sure to set and export:

```
export HTTP_PROXY=${http_proxy}
export HTTPS_PROXY=${https_proxy}

```

If you are running a local container registry, typically using Docker on port TCP:5000, you don't need to set any other variables at this point.

Otherwise, the `EMCODOCKERREPO` variables needs to be set:

```
export EMCODOCKERREPO=${container_registry_url}/
```

Note that the value for the container registry URL must end with a `/`.

Or you may set the variable in `config/config.txt`, where the default value of `EMCODOCKERREPO` is already defined.

### Update the base container images, if needed

Additional external dependencies for EMCO, and other key parameters, are also captured in the configuration file `config/config.txt`.

The configuration file specifies the following important parameters:
  * `GO_VERSION`: The version of the Go programming language to use.
  * `HELM_VERSION`: The version of the Helm package manager to use.
  * `BUILD_BASE_IMAGE_NAME`: The name of the base image used for Helm
  * `BUILD_BASE_IMAGE_VERSION`: The version of the base image used for Helm
  * `SERVICE_BASE_IMAGE_NAME`: The name of the base image for each of the EMCO microservice images.
  * `SERVICE_BASE_IMAGE_VERSION`: The name and version of the base image used for each of the EMCO microservices' images.
  * `EMCODOCKERREPO`: The container registry URL (must end with a `/`) where newly-built EMCO container images should be pushed to (default: `localhost:5000/`).
  * `MAINDOCKERREPO`: The container registry URL (must end with a `/`), where existing container images should be pulled from (default is `""` which, which results in Docker Hub).

See the file `config/config.txt` for the current default values of all the parameters above.

### Create the build base image in the EMCODOCKERREPO registry

Run the following to create the final build container image and populate that
in the `EMCODOCKERREPO` registry.

```
make build-base
```

This will build the image that can then be used for Helm.

### Registries and Images
The image names outlined in `config/config.txt` will be searched in `MAINDOCKERREPO`.
`BUILD_BASE_IMAGE_NAME` which defaults to `emco-build-base` must be first built and pushed to `EMCODOCKERREPO`. Additionally, the `mongodb` and `etcd` images will also be pulled from `MAINDOCKERREPO`.

Both `EMCODOCKERREPO` and `MAINDOCKERREPO` can be set to the same exat container registry, in which case both pulling and pushing will always be done in the same registry.

Here are some of the images and versions that have been validated as of this writing:
  1.	alpine:3.12 (this is the default `SERVICE_BASE_IMAGE_*`, the base for EMCO microservice images)
  2.	golang:1.17.7 (for building EMCO Go components)
  2.	emco-build-base:1.3 (built after `make build-base`, builds the base image used for Helm - its base being alpine:3.12 like other EMCO microservice images)
  3.	mongo:4.x
  4.	etcd:3.x

### Deploy EMCO locally
You can build and deploy the EMCO components in your local environment (and
use them to deploy your workload in a set of remote Kubernetes clusters).

This is done in two stages:

 * Build the EMCO components:
   ```make all```
   This spawns a build container that generates the needed EMCO binaries and
   container images.
 * Deploy EMCO components locally:
   ```docker-compose up```
   using `deployments/docker/docker-compose.yml`. This spawns a set of
   containers, each running one EMCO component.

See [this tutorial](docs/user/Tutorial_Local_Install.md) for further details.

### Deploy EMCO in a Kubernetes cluster
Alternatively, you can build EMCO locally and deploy EMCO components in a
Kubernetes cluster using Helm charts (and use them to deploy your workload in
another set of Kubernetes clusters).

This requires the locally built container images to be pushed to the
container registry `EMCODOCKERREPO` with the appropriate tag, so that the
Kubernetes cluster can pull images from that container registry. The tag can
be set based on whether it is a developer/test build or a release build.

Do the following steps:

 * Set up the environment:

   * For development/testing:
     ```export BUILD_CAUSE=DEV_TEST```
     This sets the image tags to the form `${USER}-latest`.

   * For release:
     ```
        export BUILD_CAUSE=RELEASE
        export EMCOSRV_RELEASE_TAG=emco-${release_number}tag
     ```
     This sets the image tags to the specified tag. Note that if you set
     `BUILD_CAUSE=RELEASE` but do not set `EMCOSRV_RELEASE_TAG`, the image tags
     will be set to any tag defined on the git `HEAD` of the current git
     branch. If no git tag is defined on the HEAD, the build will fail.

 * Set up the Helm charts: Be sure to reference those image names and tags in
   your Helm charts.

 * Build and deploy EMCO:
   ```make deploy```

### Deploy sample test cases with EMCO
See [this Readme](examples/single-cluster/Readme.md) on how to setup an environment and run a few test cases with EMCO.

See [this tutorial](docs/user/install/Tutorial_Helm.md) for further details.
