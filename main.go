package main

import (
	cli2 "github.com/urfave/cli"
	"github.com/way365/bazo-client/cli"
	"github.com/way365/bazo-client/cstorage"
	"github.com/way365/bazo-client/network"
	"github.com/way365/bazo-client/services"
	"github.com/way365/bazo-client/util"
	"github.com/way365/bazo-miner/p2p"
	"os"
)

func main() {
	p2p.InitLogging()
	services.InitLogging()
	logger := util.InitLogger()
	util.Config = util.LoadConfiguration()

	network.Init()
	cstorage.Init("client.db")

	app := cli2.NewApp()

	app.Name = "bazo-client"
	app.Usage = "the command line interface for interacting with the Bazo blockchain implemented in Go."
	app.Version = "1.0.0"
	app.Commands = []cli2.Command{
		cli.GetAccountCommand(logger),
		cli.GetFundsCommand(logger),
		cli.GetNetworkCommand(logger),
		cli.GetRestCommand(),
		cli.GetStakingCommand(logger),
		cli.GetUpdateTxCommand(logger),
	}

	err := app.Run(os.Args)
	if err != nil {
		logger.Fatal(err)
	}
}
