# syntax=docker/dockerfile:1.7

ARG GO_VERSION=1.23

FROM golang:${GO_VERSION}-bookworm AS toolchain

ENV CGO_ENABLED=1 \
    GOWORK=off \
    XLIB_CONTEXT=docker_toolchain

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
      bash \
      ca-certificates \
      curl \
      git \
      jq \
      make \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /workspace

CMD ["bash"]
