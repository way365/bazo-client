package cli

import (
	"github.com/julwil/bazo-client/args"
	"github.com/julwil/bazo-client/services"
	"github.com/urfave/cli"
	"log"
)

var (
	headerFlag = cli.IntFlag{
		Name:  "header",
		Usage: "Header flag",
		Value: 0,
	}

	feeFlag = cli.Uint64Flag{
		Name:  "fee",
		Usage: "specify the Fee",
		Value: 1,
	}

	rootkeyFlag = cli.StringFlag{
		Name:  "rootwallet",
		Usage: "load root's public private key from `FILE`",
	}
)

func GetAccountCommand(logger *log.Logger) cli.Command {
	return cli.Command{
		Name:  "account",
		Usage: "account management",
		Subcommands: []cli.Command{
			getCheckAccountCommand(logger),
			getCreateAccountCommand(logger),
			getAddAccountCommand(logger),
		},
	}
}

func getCheckAccountCommand(logger *log.Logger) cli.Command {
	return cli.Command{
		Name:  "check",
		Usage: "check account state",
		Action: func(c *cli.Context) error {
			args := &args.CheckAccountArgs{
				Address: c.String("address"),
				Wallet:  c.String("wallet"),
			}

			return services.CheckAccount(args, logger)
		},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "address",
				Usage: "the account's 128 byte Address",
			},
			cli.StringFlag{
				Name:  "wallet",
				Usage: "load the account's 128 byte Address from `FILE`",
				Value: "wallet.txt",
			},
		},
	}
}

func getCreateAccountCommand(logger *log.Logger) cli.Command {
	return cli.Command{
		Name:  "create",
		Usage: "create a new account and add it to the network",
		Action: func(c *cli.Context) error {
			args := &args.CreateAccountArgs{
				Header:     c.Int("header"),
				Fee:        c.Uint64("fee"),
				RootWallet: c.String("rootwallet"),
				Wallet:     c.String("wallet"),
				ChParams:   c.String("chparams"),
				Data:       c.String("data"),
			}

			_, err := services.PrepareSignSubmitCreateAccTx(args, logger)

			return err
		},
		Flags: []cli.Flag{
			headerFlag,
			feeFlag,
			rootkeyFlag,
			cli.StringFlag{
				Name:  "wallet",
				Usage: "save new account's public private key to `FILE`",
			},
			cli.StringFlag{
				Name:  "chparams",
				Usage: "save new chameleon hash parameters to `FILE`",
			},
			cli.StringFlag{
				Name:  "data",
				Usage: "Data field to add a message to the tx",
				Value: "",
			},
		},
	}
}

func getAddAccountCommand(logger *log.Logger) cli.Command {
	return cli.Command{
		Name:  "add",
		Usage: "add an existing account",
		Action: func(c *cli.Context) error {
			args := &args.AddAccountArgs{
				Header:     c.Int("header"),
				Fee:        c.Uint64("fee"),
				RootWallet: c.String("rootwallet"),
				Address:    c.String("address"),
				ChParams:   c.String("chparams"),
			}

			return services.AddAccount(args, logger)
		},
		Flags: []cli.Flag{
			headerFlag,
			feeFlag,
			rootkeyFlag,
			cli.StringFlag{
				Name:  "address",
				Usage: "the account's Address",
			},
			cli.StringFlag{
				Name:  "chparams",
				Usage: "save new chameleon hash parameters to `FILE`",
			},
		},
	}
}
