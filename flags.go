package main

import (
	"github.com/urfave/cli/v2"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:        "nginx.configuration",
		Usage:       "The configuration for the ingress operation of Nginx.",
		Required:    true,
		EnvVars:     []string{"NGINX_INGRESS"},
		Destination: &Pipe.Nginx.Configuration,
	},
}
