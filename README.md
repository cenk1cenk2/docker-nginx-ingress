# docker-nginx-ingress

[![pipeline status](https://gitlab.kilic.dev/docker/nginx-ingress/badges/main/pipeline.svg)](https://gitlab.kilic.dev/docker/nginx-ingress/-/commits/main) [![Docker Pulls](https://img.shields.io/docker/pulls/cenk1cenk2/nginx-ingress)](https://hub.docker.com/repository/docker/cenk1cenk2/nginx-ingress) [![Docker Image Size (latest by date)](https://img.shields.io/docker/image-size/cenk1cenk2/nginx-ingress)](https://hub.docker.com/repository/docker/cenk1cenk2/nginx-ingress) [![Docker Image Version (latest by date)](https://img.shields.io/docker/v/cenk1cenk2/nginx-ingress)](https://hub.docker.com/repository/docker/cenk1cenk2/nginx-ingress) [![GitHub last commit](https://img.shields.io/github/last-commit/cenk1cenk2/docker-nginx-ingress)](https://github.com/cenk1cenk2/docker-nginx-ingress)

# Description

Nginx ingress controller for docker-compose stacks, where it takes in the given containers as a JSON file and load-balances them through TCP/UDP stream functionality of the Nginx.

README is currently missing!

# Configuration

```json
[
  {
    "server": {
      "listen": "string",
      "options": {
        "key": "value"
      }
    },
    "upstream": {
      "servers": ["string"],
      "options": {
        "key": "value"
      }
    }
  }
]
```
