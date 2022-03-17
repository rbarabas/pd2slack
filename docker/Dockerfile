# Docker file for pd2slack service
FROM alpine:3.15 AS build

ARG PROJECT_NAME=pd2slack
ARG USER=${PROJECT_NAME}
ARG BRANCH=master
ARG VERSION=1.0.3
ARG ARCH=amd64
ARG BUILD_DIR=/build
ARG GO_DIR=${BUILD_DIR}/go
ARG SRC_DIR=${BUILD_DIR}/go/src
ARG PROJECT_DIR=${BUILD_DIR}/go/src/${PROJECT_NAME}
ARG BINARY=${PROJECT_DIR}/{PROJECT_NAME}
ARG REPO=https://github.com/rbarabas/pd2slack.git

# Root tasks
RUN rm -rf ${BUILD_DIR}
RUN mkdir ${BUILD_DIR}

RUN apk add --no-cache \
    git \
    go  \
    openssh

RUN adduser -D ${USER}
RUN chown ${USER} ${BUILD_DIR}
USER ${USER}

# Conventional user tasks
RUN mkdir ${GO_DIR}
RUN mkdir ${SRC_DIR}
WORKDIR ${BUILD_DIR}

WORKDIR ${SRC_DIR}
RUN git clone ${REPO}
WORKDIR ${PROJECT_DIR}
RUN git checkout ${BRANCH}
ENV GOPATH ${GO_DIR}
RUN go build pd2slack

FROM alpine:3.15

ARG PROJECT_NAME=pd2slack
ARG USER=${PROJECT_NAME}
ARG TARGET_DIR=/${PROJECT_NAME}
ARG BUILD_DIR=/build
ARG GO_DIR=${BUILD_DIR}/go
ARG SRC_DIR=${BUILD_DIR}/go/src

RUN mkdir ${TARGET_DIR}

RUN adduser -D ${USER}
USER ${USER}

COPY --from=build ${SRC_DIR}/pd2slack ${TARGET_DIR}
COPY --from=build ${SRC_DIR}/pd2slack/pd2slack.env ${TARGET_DIR}/pd2slack.env
COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/bin/sh", "/entrypoint.sh"]

