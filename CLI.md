# nginx-ingress

Ingress controller for applications in docker-compose stacks to do load balancing with Nginx.

`nginx-ingress [GLOBAL FLAGS] command [COMMAND FLAGS] [ARGUMENTS...]`

## Global Flags

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
|`$DEBUG` | Enable debugging for the application. | `Bool` | `false` | false |
|`$LOG_LEVEL` | Define the log level for the application.  | `String`<br/>enum(&#34;PANIC&#34;, &#34;FATAL&#34;, &#34;WARNING&#34;, &#34;INFO&#34;, &#34;DEBUG&#34;, &#34;TRACE&#34;) | `false` | &#34;info&#34; |
|`$NGINX_INGRESS` | The configuration for the ingress operation of Nginx. | `String` | `true` |  |

## Commands

### `help` , `h`

`Shows a list of commands or help for one command`
