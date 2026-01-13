package pipe

import (
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:     "nginx.configuration",
		Usage:    "The configuration for the ingress operation of Nginx. json({ server: struct { listen: string, options: map[string]string }, upstream: struct { servers: []string, options: map[string]string } })",
		Required: true,
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("NGINX_INGRESS"),
		),
		Validator: func(v string) error {
			if err := json.Unmarshal([]byte(v), &P.Nginx.Configuration); err != nil {
				return fmt.Errorf("Can not unmarshal configuration: %w", err)
			}

			return nil

		},
	},
}
