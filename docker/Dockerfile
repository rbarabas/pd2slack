# Docker file for pd2slack service
FROM alpine:3.15 AS build

ARG BRANCH=poc
ARG VERSION=1.0.0
ARG ARCH=amd64

RUN mkdir /build
WORKDIR /build

RUN apk add --no-cache \
    git \
    go

RUN adduser -D pd2slack
RUN chown pd2slack /build
USER pd2slack

RUN git clone https://github.com/BoltApp/pd2slack.git
WORKDIR /build/pd2slack
RUN git checkout poc
RUN go build pd2slack

FROM alpine:3.15

COPY --from=build /build/pd2slack /usr/bin

RUN adduser -D pd2slack
USER pd2slack

COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/bin/sh", "/entrypoint.sh"]

