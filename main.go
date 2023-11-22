package main

import (
	"github.com/urfave/cli/v2"

	"gitlab.kilic.dev/docker/nginx-ingress/pipe"
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

func main() {
	NewPlumber(
		func(p *Plumber) *cli.App {
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
						pipe.New(p).SetCliContext(c).Job(),
					)
				},
			}
		}).
		SetDocumentationOptions(DocumentationOptions{
			MarkdownOutputFile: "CLI.md",
			MarkdownBehead:     0,
			ExcludeFlags:       true,
			ExcludeHelpCommand: true,
		}).
		Run()
}
