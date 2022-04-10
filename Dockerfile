FROM nginx:alpine

COPY ./dist/pipe /usr/bin/pipe

RUN chmod +x /usr/bin/pipe && \
  # smoke test
  pipe --help

COPY ./.docker/hostfs /

WORKDIR /etc/nginx

ENTRYPOINT [ "pipe" ]
