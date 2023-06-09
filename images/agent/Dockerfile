# Copyright 2023 Authors of kdoctor-io
# SPDX-License-Identifier: Apache-2.0

#====modify====
ARG BASE_IMAGE=ghcr.io/kdoctor-io/kdoctor-agent-base:bd4eba101d4fb9c6fa5c90c704fff5c6446bda70
ARG GOLANG_IMAGE=docker.io/library/golang:1.20

#======= build bin ==========
# use BUILDPLATFORM to cross build
FROM --platform=${BUILDPLATFORM} ${GOLANG_IMAGE} as builder

# docker build var
ARG TARGETOS
ARG TARGETARCH

# custom var
ARG RACE
ARG NOSTRIP
ARG NOOPT

COPY . /src
WORKDIR /src
RUN  make GOARCH=${TARGETARCH}   \
        RACE=${RACE} NOSTRIP=${NOSTRIP} NOOPT=${NOOPT} \
        DESTDIR_BIN=/tmp/install/${TARGETOS}/${TARGETARCH}/bin \
        build_agent_bin



#====== release image =======

FROM --platform=${TARGETPLATFORM} ${BASE_IMAGE}

LABEL maintainer="maintainer@kdoctor-io"

# TARGETOS is an automatic platform ARG enabled by Docker BuildKit.
ARG TARGETOS
# TARGETARCH is an automatic platform ARG enabled by Docker BuildKit.
ARG TARGETARCH

ARG GIT_COMMIT_VERSION
ENV ENV_GIT_COMMIT_VERSION=${GIT_COMMIT_VERSION}
ARG GIT_COMMIT_TIME
ENV ENV_GIT_COMMIT_TIMESTAMP=${GIT_COMMIT_TIME}
ARG VERSION
ENV ENV_VERSION=${VERSION}

COPY --from=builder /tmp/install/${TARGETOS}/${TARGETARCH}/bin/*   /usr/bin/

CMD ["/usr/bin/agent"]
