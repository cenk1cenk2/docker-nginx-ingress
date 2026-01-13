# syntax=docker/dockerfile-upstream:master-labs
FROM nginx:alpine

RUN \
  apk add --no-cache tini

ARG BUILDOS
ARG BUILDARCH

COPY --chmod=777 ./dist/pipe-${BUILDOS}-${BUILDARCH} /usr/bin/pipe

RUN \
  # smoke test
  pipe --help

WORKDIR /etc/nginx

ENTRYPOINT [ "tini", "pipe" ]
