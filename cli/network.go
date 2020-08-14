package cli

import (
	"errors"
	"github.com/julwil/bazo-client/args"
	"github.com/julwil/bazo-client/services"
	"github.com/urfave/cli"
	"log"
)

func GetNetworkCommand(logger *log.Logger) cli.Command {
	options := []args.ConfigOption{
		{Id: 1, Name: "setBlockSize", Usage: "set the size of blocks (in bytes)"},
		{Id: 2, Name: "setDifficultyInterval", Usage: "set the difficulty interval (in number of blocks)"},
		{Id: 3, Name: "setMinimumFee", Usage: "set the minimum Fee (in Bazo coins)"},
		{Id: 4, Name: "setBlockInterval", Usage: "set the block interval (in seconds)"},
		{Id: 5, Name: "setBlockReward", Usage: "set the block reward (in Bazo coins)"},
	}

	command := cli.Command{
		Name:  "network",
		Usage: "configure the network",
		Action: func(c *cli.Context) error {
			optionsSetByUser := 0
			for _, option := range options {
				if !c.IsSet(option.Name) {
					continue
				}

				optionsSetByUser++

				args := &args.NetworkArgs{
					Header:     c.Int("Header"),
					Fee:        c.Uint64("Fee"),
					TootWallet: c.String("rootwallet"),
					OptionId:   option.Id,
					Payload:    c.Uint64(option.Name),
					TxCount:    c.Int("TxCount"),
				}

				err := services.ConfigureNetwork(args, logger)
				if err != nil {
					return err
				}
			}

			if optionsSetByUser == 0 {
				return errors.New("specify at least one configuration option")
			}

			return nil
		},
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "Header",
				Usage: "Header flag",
				Value: 0,
			},
			cli.Uint64Flag{
				Name:  "Fee",
				Usage: "specify the Fee",
				Value: 1,
			},
			cli.IntFlag{
				Name:  "TxCount",
				Usage: "the sender's current transaction counter",
			},
			cli.StringFlag{
				Name:  "rootwallet",
				Usage: "load root's public key from `FILE`",
			},
		},
	}

	for _, option := range options {
		flag := cli.Uint64Flag{Name: option.Name, Usage: option.Usage}
		command.Flags = append(command.Flags, flag)
	}

	return command
}
