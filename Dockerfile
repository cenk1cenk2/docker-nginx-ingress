FROM nginx:alpine

COPY ./dist/pipe /usr/bin/pipe

RUN \
  apk add --no-cache tini && \
  chmod +x /usr/bin/pipe && \
  # smoke test
  pipe --help

COPY ./.docker/hostfs /

WORKDIR /etc/nginx

ENTRYPOINT [ "tini", "pipe" ]
