# syntax=docker/dockerfile-upstream:master-labs
FROM nginx:alpine

RUN \
  apk add --no-cache tini

COPY --chmod=777 ./dist/pipe-${BUILDOS}-${BUILDARCH} /usr/bin/pipe

RUN \
  # smoke test
  pipe --help

COPY ./.docker/hostfs /

WORKDIR /etc/nginx

ENTRYPOINT [ "tini", "pipe" ]
