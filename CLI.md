# nginx-ingress

Ingress controller for applications in docker-compose stacks to do load balancing with Nginx.

`nginx-ingress [FLAGS]`

## Flags

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$NGINX_INGRESS` | The configuration for the ingress operation of Nginx. | `string`<br/>`json({ server: struct { listen: string, options: map[string]string }, upstream: struct { servers: []string, options: map[string]string } })` | `true` | <code></code> |

**CLI**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application. | `string`<br/>`enum("panic", "fatal", "warn", "info", "debug", "trace")` | `false` | <code>"info"</code> |
| `$ENV_FILE` | Environment files to inject. | `string[]` | `false` | <code></code> |
