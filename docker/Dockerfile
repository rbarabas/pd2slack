# Docker file for pd2slack service
FROM alpine:3.15 AS build

ARG BRANCH=poc
ARG VERSION=1.0.3
ARG ARCH=amd64

RUN mkdir /build
WORKDIR /build

RUN apk add --no-cache \
    git \
    go

RUN adduser -D pd2slack
RUN chown pd2slack /build
USER pd2slack

RUN rm -rf /build/pdslack
RUN git clone https://github.com/rbarabas/pd2slack.git
WORKDIR /build/pd2slack
RUN git checkout poc
RUN go build pd2slack

FROM alpine:3.15

RUN mkdir /pd2slack
COPY --from=build /build/pd2slack /pd2slack
COPY --from=build /build/pd2slack/pd2slack.env /pd2slack/pd2slack.env

RUN adduser -D pd2slack
USER pd2slack

COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/bin/sh", "/entrypoint.sh"]

