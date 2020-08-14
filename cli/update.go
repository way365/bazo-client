package cli

import (
	"github.com/julwil/bazo-client/args"
	"github.com/julwil/bazo-client/services"
	"github.com/urfave/cli"
	"log"
)

func GetUpdateTxCommand(logger *log.Logger) cli.Command {
	return cli.Command{
		Name:  "update",
		Usage: "update the data field of a specific transaction",
		Action: func(c *cli.Context) error {
			args := &args.UpdateTxArgs{
				Header:     c.Int("header"),
				Fee:        c.Uint64("fee"),
				TxToUpdate: c.String("tx-hash"),
				TxIssuer:   c.String("tx-issuer"),
				ChParams:   c.String("chparams"),
				UpdateData: c.String("update-data"),
				Data:       c.String("data"),
			}

			err := args.ValidateInput()
			if err != nil {
				return err
			}

			_, err = services.PrepareSignSubmitUpdateTx(args, logger)
			return err
		},
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "header",
				Usage: "Header flag",
				Value: 0,
			},
			cli.Uint64Flag{
				Name:  "fee",
				Usage: "specify the Fee",
				Value: 1,
			},
			cli.StringFlag{
				Name:  "tx-hash",
				Usage: "the 32-byte hash of the transaction to be upddated",
			},
			cli.StringFlag{
				Name:  "tx-issuer",
				Usage: "load the tx issuer's public key from `FILE`",
			},
			cli.StringFlag{
				Name:  "chparams",
				Usage: "load the chameleon hash parameters from `FILE`",
			},
			cli.StringFlag{
				Name:  "update-data",
				Usage: "specify the new Data that shall be updated on the tx",
			},
			cli.StringFlag{
				Name:  "data",
				Usage: "specify the Data on this tx.",
			},
		},
	}
}
