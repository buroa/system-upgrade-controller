FROM golang:1.22-alpine3.19 AS builder

ARG PKG=github.com/buroa/system-upgrade-controller
ARG VERSION=dev
ARG REVISION=dev
ARG BUILDTIME

WORKDIR /src

COPY . .

RUN go build -ldflags "-s -w -X ${PKG}/pkg/version.Version=${VERSION} -X ${PKG}/pkg/version.GitCommit=${REVISION}" -o bin/system-upgrade-controller

FROM scratch

LABEL org.opencontainers.image.source = "https://github.com/buroa/system-upgrade-controller"

COPY --from=builder /src/bin/system-upgrade-controller /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/system-upgrade-controller"]
