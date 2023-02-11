package pipe

import "embed"

var (
	//go:embed templates
	Templates embed.FS
)

const (
	TEMPLATE_FOLDER_SERVERS         = "servers"
	TEMPLATE_FOLDER_UPSTREAMS       = "upstreams"
	NGINX_ROOT_CONFIGURATION_FOLDER = "/etc/nginx"
	NGINX_CONFIGURATION             = "nginx.conf"
)
