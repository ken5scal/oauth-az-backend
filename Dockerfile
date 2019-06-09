ARG GO_VERSION=1.12
ARG ALPINE_VERSION=3.9
ARG PORT=8080
ARG ENV="debug"

# This file is inspired from following articles
# https://medium.com/@pierreprinetti/the-go-1-11-dockerfile-a3218319d191
# https://qiita.com/takasp/items/c6288d4836e79801bb19#dockerfile-1
# https://qiita.com/theoden9014/items/92c598d6662bd6c6b194

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

# By separating code copying from dependency-related files,
# build steps can be cached.
COPY . .
ENV CGO_ENABLED=0 GOOS=linux
RUN go build -ldflags="-w -s" -o /go/bin/app

## Runtime Library
FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/bin/app /go/bin/app

ENV PORT ${PORT}
ENV ENVIRONMENT ${ENV}
EXPOSE ${PORT}

USER app-go
ENTRYPOINT ["/go/bin/app"]