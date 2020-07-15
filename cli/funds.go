package cli

import (
	"github.com/julwil/bazo-client/args"
	"github.com/julwil/bazo-client/client"
	"github.com/urfave/cli"
	"log"
)

func GetFundsCommand(logger *log.Logger) cli.Command {
	return cli.Command{
		Name:  "funds",
		Usage: "send funds from one account to another",
		Action: func(c *cli.Context) error {
			args := &args.FundsArgs{
				Header:         c.Int("header"),
				FromWalletFile: c.String("from"),
				ToWalletFile:   c.String("to"),
				ToAddress:      c.String("toAddress"),
				MultisigFile:   c.String("multisig"),
				ChParamsFile:   c.String("chparams"),
				Amount:         c.Uint64("amount"),
				Fee:            c.Uint64("fee"),
				TxCount:        c.Int("txcount"),
				Data:           c.String("data"),
			}

			err := args.ValidateInput()
			if err != nil {
				return err
			}

			_,err = client.SendFunds(args, logger)

			return err
		},
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "header",
				Usage: "Header flag",
				Value: 0,
			},
			cli.StringFlag{
				Name:  "from",
				Usage: "load the sender's private key from `FILE`",
			},
			cli.StringFlag{
				Name:  "to",
				Usage: "load the recipient's public key from `FILE`",
			},
			cli.StringFlag{
				Name:  "ToAddress",
				Usage: "the recipient's 128 byze public address",
			},
			cli.Uint64Flag{
				Name:  "amount",
				Usage: "specify the Amount to send",
			},
			cli.Uint64Flag{
				Name:  "fee",
				Usage: "specify the Fee",
				Value: 1,
			},
			cli.IntFlag{
				Name:  "txcount",
				Usage: "the sender's current transaction counter",
			},
			cli.StringFlag{
				Name:  "multisig",
				Usage: "load multi-signature server’s private key from `FILE`",
			},
			cli.StringFlag{
				Name:  "chparams",
				Usage: "load the chameleon hash parameters from `FILE`",
			},
			cli.StringFlag{
				Name:  "data",
				Usage: "Data field to add a message to the tx",
			},
		},
	}
}
