package cli

import (
	"github.com/julwil/bazo-client/args"
	"github.com/julwil/bazo-client/client"
	"github.com/urfave/cli"
	"log"
)

func GetStakingCommand(logger *log.Logger) cli.Command {
	headerFlag := cli.IntFlag{
		Name:  "header",
		Usage: "Header flag",
		Value: 0,
	}

	feeFlag := cli.Uint64Flag{
		Name:  "fee",
		Usage: "specify the Fee",
		Value: 1,
	}

	walletFlag := cli.StringFlag{
		Name:  "wallet, w",
		Usage: "load validator's public key from `FILE`",
		Value: "wallet.txt",
	}

	return cli.Command{
		Name:  "staking",
		Usage: "enable or disable staking",
		Subcommands: []cli.Command{
			{
				Name:  "enable",
				Usage: "join the pool of validators",
				Action: func(c *cli.Context) error {
					args := args.ParseStakingArgs(c)
					args.StakingValue = true
					return client.ToggleStaking(args, logger)
				},
				Flags: []cli.Flag{
					headerFlag,
					feeFlag,
					walletFlag,
					cli.StringFlag{
						Name:  "commitment",
						Usage: "load valiadator's Commitment key from `FILE`",
						Value: "Commitment.txt",
					},
				},
			},
			{
				Name:  "disable",
				Usage: "leave the pool of validators",
				Action: func(c *cli.Context) error {
					args := args.ParseStakingArgs(c)
					args.StakingValue = false
					return client.ToggleStaking(args, logger)
				},
				Flags: []cli.Flag{
					headerFlag,
					feeFlag,
					walletFlag,
				},
			},
		},
	}
}
