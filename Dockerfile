# syntax=docker/dockerfile:1.7

ARG GO_VERSION=1.23
ARG GO_BASE_IMAGE=golang:${GO_VERSION}-bookworm
ARG GO_BASE_IMAGE_DIGEST=sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
ARG GOLANGCI_LINT_VERSION=v2.1.6
ARG GOVULNCHECK_VERSION=v1.3.0

FROM ${GO_BASE_IMAGE} AS toolchain

ARG GO_BASE_IMAGE
ARG GO_BASE_IMAGE_DIGEST
ARG GOLANGCI_LINT_VERSION
ARG GOVULNCHECK_VERSION

LABEL org.opencontainers.image.base.name="${GO_BASE_IMAGE}"
LABEL org.opencontainers.image.base.digest="${GO_BASE_IMAGE_DIGEST}"

ENV CGO_ENABLED=1 \
    GOWORK=off \
    XLIB_CONTEXT=docker_toolchain \
    GOPATH=/go \
    GOCACHE=/root/.cache/go-build \
    GOLANGCI_LINT_VERSION=${GOLANGCI_LINT_VERSION} \
    GOVULNCHECK_VERSION=${GOVULNCHECK_VERSION}

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
      bash \
      ca-certificates \
      curl \
      git \
      jq \
      make \
      openssh-client \
      tar \
      xz-utils \
    && rm -rf /var/lib/apt/lists/*

# golangci-lint: 使用预编译二进制，避免 go install 下载大量传递依赖导致超时
RUN curl -sSfL https://github.com/golangci/golangci-lint/releases/download/${GOLANGCI_LINT_VERSION}/golangci-lint-${GOLANGCI_LINT_VERSION#v}-linux-amd64.tar.gz \
    | tar -xz --strip-components=1 -C /usr/local/bin golangci-lint-${GOLANGCI_LINT_VERSION#v}-linux-amd64/golangci-lint

# govulncheck: v1.2.0+ 要求 Go >= 1.25，Go 1.23 环境使用 v1.1.4（兼容 go 1.22+）
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go install golang.org/x/vuln/cmd/govulncheck@v1.1.4

WORKDIR /workspace

CMD ["bash"]

FROM toolchain AS dev

ARG USERNAME=xlib
ARG USER_UID=1000
ARG USER_GID=1000

RUN groupadd --gid ${USER_GID} ${USERNAME} \
    && useradd --uid ${USER_UID} --gid ${USER_GID} -m ${USERNAME} \
    && mkdir -p /workspace /home/${USERNAME}/.cache/go-build /go/pkg/mod \
    && chown -R ${USERNAME}:${USERNAME} /workspace /home/${USERNAME} /go

USER ${USERNAME}
ENV XLIB_CONTEXT=docker_toolchain \
    GOCACHE=/home/${USERNAME}/.cache/go-build

CMD ["bash"]

FROM toolchain AS gate

COPY . /workspace

CMD ["make", "ci"]

FROM toolchain AS build-goalcli

COPY . /workspace
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/goalcli ./cmd/goalcli

FROM scratch AS goalcli-runtime

COPY --from=build-goalcli /out/goalcli /goalcli

ENTRYPOINT ["/goalcli"]
CMD ["version"]
