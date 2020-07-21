package cli

import (
	"github.com/julwil/bazo-client/client"
	"github.com/julwil/bazo-client/rest"
	"github.com/urfave/cli"
)

func GetRestCommand() cli.Command {
	return cli.Command{
		Name:  "rest",
		Usage: "start the rest service",
		Action: func(c *cli.Context) error {
			client.Sync()
			rest.Init()
			return nil
		},
	}
}
