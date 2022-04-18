# docker-nginx-ingress

[![pipeline status](https://gitlab.kilic.dev/docker/nginx-ingress/badges/main/pipeline.svg)](https://gitlab.kilic.dev/docker/nginx-ingress/-/commits/main) [![Docker Pulls](https://img.shields.io/docker/pulls/cenk1cenk2/nginx-ingress)](https://hub.docker.com/repository/docker/cenk1cenk2/nginx-ingress) [![Docker Image Size (latest by date)](https://img.shields.io/docker/image-size/cenk1cenk2/nginx-ingress)](https://hub.docker.com/repository/docker/cenk1cenk2/nginx-ingress) [![Docker Image Version (latest by date)](https://img.shields.io/docker/v/cenk1cenk2/nginx-ingress)](https://hub.docker.com/repository/docker/cenk1cenk2/nginx-ingress) [![GitHub last commit](https://img.shields.io/github/last-commit/cenk1cenk2/docker-nginx-ingress)](https://github.com/cenk1cenk2/docker-nginx-ingress)

## Description

Nginx ingress controller for docker-compose stacks, where it takes in a JSON environment variable for the defined containers and load-balances them through streams of Nginx.

<!-- toc -->

- [Setup](#setup)
  - [Environment Variables](#environment-variables)
    - [`NGINX_INGRESS`](#nginx_ingress)

<!-- tocstop -->

## Setup

You can run this application as a `docker-compose` stack. The image is hosted as `cenk1cenk2/nginx-ingress` on DockerHub. Check out the [docker-compose](./docker-compose.yml) file for example configuration.

### Environment Variables

| Variable      | Description                                                                                                 |
| ------------- | ----------------------------------------------------------------------------------------------------------- |
| LOG_LEVEL     | Sets the log level for the container. Can take values of "fatal", "error", "warn", "info", "debug", "trace" |
| NGINX_INGRESS | JSON format configuration to setup the containers to load-balance.                                          |

#### `NGINX_INGRESS`

The environment variable `NGINX_INGRESS` is an array of objects in the JSON form to define the endpoints and load-balanced containers.

```jsonc
[
  {
    "server": {
      "listen": "string", // listen port and type for endpoint
      "options": {
        // key-value pairs of options that should be passed to "server" configuration of Nginx
      }
    },
    "upstream": {
      "servers": [
        // string slice of balanced servers
      ],
      "options": {
        // key-value pairs of options that should be passed to "upstream" configuration of Nginx
      }
    }
  }
]
```
