package main

import (
	"github.com/urfave/cli/v2"

	"gitlab.kilic.dev/docker/nginx-ingress/pipe"
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

func main() {
	p := Plumber{
		DocsFile:               "CLI.md",
		DocsExcludeFlags:       true,
		DocsExcludeHelpCommand: true,
	}

	p.New(
		func(a *Plumber) *cli.App {
			return &cli.App{
				Name:        CLI_NAME,
				Version:     VERSION,
				Usage:       DESCRIPTION,
				Description: DESCRIPTION,
				Flags:       pipe.Flags,
				Before: func(ctx *cli.Context) error {
					p.EnableTerminator()

					return nil
				},
				Action: func(c *cli.Context) error {
					return pipe.TL.RunJobs(
						pipe.New(a).SetCliContext(c).Job(),
					)
				},
			}
		}).Run()
}
