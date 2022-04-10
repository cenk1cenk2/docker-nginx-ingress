package main

import (
	"github.com/urfave/cli/v2"

	utils "github.com/cenk1cenk2/ci-cd-pipes/utils"
)

func main() {
	utils.CliCreate(
		&cli.App{
			Name:        CLI_NAME,
			Version:     VERSION,
			Usage:       DESCRIPTION,
			Description: DESCRIPTION,
			Action: func(c *cli.Context) error {
				utils.CliGreet(c)

				err := cli.ShowAppHelp(c)

				if err != nil {
					return err
				}

				utils.Log.Fatalln("Need a subcommand to run!")

				return nil
			},
		},
	)
}
