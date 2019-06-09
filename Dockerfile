# This file is inspired from following articles
# https://medium.com/@pierreprinetti/the-go-1-11-dockerfile-a3218319d191
# https://qiita.com/takasp/items/c6288d4836e79801bb19#dockerfile-1
# https://qiita.com/theoden9014/items/92c598d6662bd6c6b194
# https://qiita.com/theoden9014/items/92c598d6662bd6c6b194

ARG GO_VERSION=1.12
ARG ALPINE_VERSION=3.9

# -------------------------------------------------------
# ------------- Stage where building go app -------------
# -------------------------------------------------------
# docker build -f Dockerfile -t oauth-az-back-build --target builder .
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder
MAINTAINER Kengo Suzuki <kengoscal@gmail.com>

RUN addgroup -S -g 50001 app && \
    adduser -D -S -G app -u 50001 app-go

WORKDIR $GOPATH/src/project/

# Pre requisites running `go mod download`
RUN apk update && apk --no-cache add \
    git \
    ca-certificates && \
    rm -rf /var/cache/apk/*
ENV GO111MODULE=on

# Download Dependencies
COPY go.mod go.sum ./
RUN go mod download

# Reading Dependency file
COPY config.toml /etc/oauth-az/

# By separating code copying from dependency-related files,
# build steps can be cached.
COPY . .
ENV CGO_ENABLED=0 GOOS=linux
RUN go build -ldflags="-w -s" -o /go/bin/app

# ------------------------------------------------------
# ------------- Stage where running go app -------------
# ------------------------------------------------------
# % docker build -f Dockerfile --build-arg ENV=debug -t oauth-az-back-dev .
FROM scratch

COPY --from=builder /etc/group /etc/passwd /etc/
COPY --from=builder /etc/oauth-az /etc/oauth-az
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/bin/app /go/bin/app

USER app-go
ENTRYPOINT ["/go/bin/app", "/etc/oauth-az/config.toml"]