FROM golang:alpine AS build_base

RUN apk add --no-cache git curl make build-base upx

RUN mkdir -p /app

WORKDIR /app

ENV GOPROXY="https://proxy.golang.org,direct"

COPY go.* ./

RUN go mod download


FROM build_base as builder

COPY . .

RUN make docker_install

RUN strip /go/bin/keybase-docker


FROM alpine

WORKDIR /

COPY --from=builder /go/bin/keybase-docker /keybase-docker

ENTRYPOINT ["./keybase-docker"]

