# nginx-ingress

Ingress controller for applications in docker-compose stacks to do load balancing with Nginx.

`nginx-ingress [FLAGS]`

## Flags

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NGINX_INGRESS` | The configuration for the ingress operation of Nginx.  | `String`<br/>`json({ server: struct { listen: string, options: map[string]string }, upstream: struct { servers: []string, options: map[string]string } })` | `true` |  |

### CLI

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application.  | `String`<br/>`enum("PANIC", "FATAL", "WARNING", "INFO", "DEBUG", "TRACE")` | `false` | info |
| `$ENV_FILE` | Environment files to inject. | `StringSlice` | `false` |  |
