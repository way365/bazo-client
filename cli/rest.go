package cli

import (
	"github.com/julwil/bazo-client/http"
	"github.com/julwil/bazo-client/services"
	"github.com/urfave/cli"
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
