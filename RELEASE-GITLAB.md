# Release steps for emco-base on GitLab

This document outlines the release process for emco-base on GitLab.
The version as of this writing is `22.06`. Please modify to the correct version according to the release undertaken.

## Tagging the source code

    git tag v22.06
    git push --tags

## Pushing container images

    cd emco-base # this is the git repo dir
    docker login registry.gitlab.com

    export BUILD_CAUSE=RELEASE
    export EMCOSRV_RELEASE_TAG=v22.06
    export EMCODOCKERREPO=registry.gitlab.com/project-emco/core/emco-base/
    make deploy

## Updating Helm charts

Usually, there is no need to update the public Helm charts hosted on the GitLab Package Registry, since they are designed to be compatible with any version of EMCO hosted on the GitLab Container Registry.
However, if indeed there was a bump the any Helm chart/package `version`, including when new microservices get added (which require a bump to the `version` of the encompassing packages, i.e. `emco-services` followed by `emco`), then here is how to generate and upload new packages:

    PROJECT_ACCESS_TOKEN="REPLACE_WITH_GITLAB_PROJECT_ACCESS_TOKEN"
    cd emco-base/deployments/helm/emcoBase
    make clean
    make all
    curl --request POST --form chart=@dist/packages/emco-db-1.0.0.tgz --user emco-base:$PROJECT_ACCESS_TOKEN https://gitlab.com/api/v4/projects/29353813/packages/helm/api/stable/charts
    curl --request POST --form chart=@dist/packages/emco-tools-1.0.0.tgz --user emco-base:$PROJECT_ACCESS_TOKEN https://gitlab.com/api/v4/projects/29353813/packages/helm/api/stable/charts
    curl --request POST --form chart=@dist/packages/emco-services-1.0.1.tgz --user emco-base:$PROJECT_ACCESS_TOKEN https://gitlab.com/api/v4/projects/29353813/packages/helm/api/stable/charts
    curl --request POST --form chart=@dist/packages/emco-1.0.1.tgz --user emco-base:$PROJECT_ACCESS_TOKEN https://gitlab.com/api/v4/projects/29353813/packages/helm/api/stable/charts
    cd ..
    helm package monitor
    curl --request POST --form chart=@monitor-1.0.0.tgz --user emco-base:$PROJECT_ACCESS_TOKEN https://gitlab.com/api/v4/projects/29353813/packages/helm/api/stable/charts

The commands above are examples for the Helm charts that were released as of EMCO 22.06.

Make sure to replace set `$PROJECT_ACCESS_TOKEN` to a [GitLab project access token](https://docs.gitlab.com/ee/user/project/settings/project_access_tokens.html) for EMCO.

Currently only the `stable` channel of the package registry is being used, but we may decide to have multiple channels in the future once we need to support charts of different versions that shouldn't clash with each other (for example to prevent a stable version from upgrading to another stable version, and instead only upgrade to point releases).

## Creating release page and tarballs

Go to the [GitLab Releases page](https://gitlab.com/project-emco/core/emco-base/-/releases), and click *New release*.
Choose the tag name `v22.06` created earlier. Name the release title simply as `22.06`. Choose the corresponding milestone, i.e. 22.06. Finally, write the release notes, add any extra release assets (see subsection below), and click *Create release*.

### Extra release assets

The emcoctl tool should be added as a release asset with the name `emcoctl-${GOOS}-${GOARCH}`, for example `emcoctl-linux-amd64`. The binary and checksum must first be uploaded to the generic package repository, and once uploaded emcoctl will be present in the package registry asset linked on the release page.

    cd bin/emcoctl
    cp emcoctl emcoctl-linux-amd64
    sha256sum emcoctl-linux-amd64 >emcoctl-linux-amd64.sha256
    PROJECT_ACCESS_TOKEN="REPLACE_WITH_GITLAB_PROJECT_ACCESS_TOKEN"
    curl --header "PRIVATE-TOKEN: $PROJECT_ACCESS_TOKEN" --upload-file emcoctl-linux-amd64 https://gitlab.com/api/v4/projects/29353813/packages/generic/emcoctl/v22.06/emcoctl-linux-amd64
    curl --header "PRIVATE-TOKEN: $PROJECT_ACCESS_TOKEN" --upload-file emcoctl-linux-amd64.sha256 https://gitlab.com/api/v4/projects/29353813/packages/generic/emcoctl/v22.06/emcoctl-linux-amd64.sha256
