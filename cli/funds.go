package cli

import (
	"github.com/urfave/cli"
	"github.com/way365/bazo-client/args"
	"github.com/way365/bazo-client/services"
	"log"
)

func GetFundsCommand(logger *log.Logger) cli.Command {
	return cli.Command{
		Name:  "funds",
		Usage: "send funds from one account to another",
		Action: func(c *cli.Context) error {
			args := &args.FundsArgs{
				Header:      c.Int("header"),
				From:        c.String("from"),
				To:          c.String("to"),
				MultiSigKey: c.String("multisig"),
				Parameters:  c.String("chparams"),
				Amount:      c.Uint64("amount"),
				Fee:         c.Uint64("fee"),
				TxCount:     c.Int("txcount"),
				Data:        c.String("data"),
			}

			err := args.ValidateInput()
			if err != nil {
				return err
			}

			_, err = services.PrepareSignSubmitFundsTx(args, logger)

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
				Usage: "load the sender's private key from `FILE` or provide the private key directly",
			},
			cli.StringFlag{
				Name:  "to",
				Usage: "load the recipient's public key from `FILE` or provide the public key directly",
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
				Usage: "load the chameleon hash parameters from `FILE` or provide them directly",
			},
			cli.StringFlag{
				Name:  "data",
				Usage: "Data field to add a message to the tx",
				Value: "",
			},
		},
	}
}
