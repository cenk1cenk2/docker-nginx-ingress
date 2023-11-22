package pipe

import (
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v2"
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

//revive:disable:line-length-limit

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:     "nginx.configuration",
		Usage:    "The configuration for the ingress operation of Nginx. json({ server: struct { listen: string, options: map[string]string }, upstream: struct { servers: []string, options: map[string]string } })",
		Required: true,
		EnvVars:  []string{"NGINX_INGRESS"},
	},
}

func ProcessFlags(tl *TaskList[Pipe]) error {
	if v := tl.CliContext.String("nginx.configuration"); v != "" {
		if err := json.Unmarshal([]byte(v), &tl.Pipe.Nginx.Configuration); err != nil {
			return fmt.Errorf("Can not unmarshal configuration: %w", err)
		}
	}

	return nil
}
