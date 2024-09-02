package cli

import (
	"github.com/urfave/cli"
	"github.com/way365/bazo-client/http"
	"github.com/way365/bazo-client/services"
)

func GetRestCommand() cli.Command {
	return cli.Command{
		Name:  "rest",
		Usage: "start the rest service",
		Action: func(c *cli.Context) error {
			services.Sync()
			http.Init()
			return nil
		},
	}
}
