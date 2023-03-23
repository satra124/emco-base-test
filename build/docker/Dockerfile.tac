# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2022 Intel Corporation

ARG MAINDOCKERREPO
ARG SERVICE_BASE_IMAGE_NAME
ARG SERVICE_BASE_IMAGE_VERSION

FROM ${MAINDOCKERREPO}${SERVICE_BASE_IMAGE_NAME}:${SERVICE_BASE_IMAGE_VERSION}

ARG EMCO_SHA
ARG EMCO_VERSION
ENV EMCO_META_EMCO_SHA=${EMCO_SHA}
ENV EMCO_META_EMCO_VERSION=${EMCO_VERSION}

WORKDIR /opt/emco/tac

RUN addgroup -S emco && adduser -S -G emco emco
RUN chown emco:emco . -R

COPY --chown=emco ./tac ./
COPY --chown=emco ./config.json ./
COPY --chown=emco ./json-schemas ./json-schemas
COPY --chown=emco ./ref-schemas ./ref-schemas

USER emco

ENTRYPOINT ["./tac"]
