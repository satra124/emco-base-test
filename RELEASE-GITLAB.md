# Release steps for emco-base on GitLab

This document outlines the release process for emco-base on GitLab.
The version as of this writing is `22.03`. Please modify to the correct version according to the release undertaken.

## Tagging the source code

    git tag v22.03
    git push --tags

## Pushing container images

    cd emco-base # this is the git repo dir
    docker login registry.gitlab.com

    export BUILD_CAUSE=RELEASE
    export EMCOSRV_RELEASE_TAG=v22.03
    export EMCODOCKERREPO=registry.gitlab.com/project-emco/core/emco-base/
    make deploy

## Creating release page and tarballs

Go to the [GitLab Releases page](https://gitlab.com/project-emco/core/emco-base/-/releases), and click *New release*.
Choose the tag name `v22.03` created earlier. Name the release title simply as `22.03`. Choose the corresponding milestone, i.e. 22.03. Finally, write the release notes, add any extra release assets, and click *Create release*.
