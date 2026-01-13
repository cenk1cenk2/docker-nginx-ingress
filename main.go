package main

import (
	"context"

	"github.com/urfave/cli/v3"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/docker/nginx-ingress/pipe"
)

func main() {
	NewPlumber(
		func(p *Plumber) *cli.Command {
			return &cli.Command{
				Name:        CLI_NAME,
				Version:     VERSION,
				Usage:       DESCRIPTION,
				Description: DESCRIPTION,
				Flags:       pipe.Flags,
				Before: func(ctx context.Context, _ *cli.Command) (context.Context, error) {
					p.EnableTerminator()

					return ctx, nil
				},
				Action: func(ctx context.Context, _ *cli.Command) error {
					return p.RunJobs(
						CombineTaskLists(
							pipe.New(p),
						),
					)
				},
			}
		}).
		SetDocumentationOptions(DocumentationOptions{
			MarkdownOutputFile: "CLI.md",
			MarkdownBehead:     0,
			ExcludeFlags:       true,
		}).
		Run()
}
