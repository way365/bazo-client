package cli

import (
	"errors"
	"github.com/julwil/bazo-miner/crypto"
	"github.com/julwil/bazo-miner/protocol"
	"github.com/urfave/cli"
	"log"
	"math/big"
)

type addAccountArgs struct {
	header             int
	fee                uint64
	rootWalletFile     string
	address            string
	chamHashParamsFile string
}

func getAddAccountCommand(logger *log.Logger) cli.Command {
	return cli.Command{
		Name:  "add",
		Usage: "add an existing account",
		Action: func(c *cli.Context) error {
			args := &addAccountArgs{
				header:             c.Int("Header"),
				fee:                c.Uint64("Fee"),
				rootWalletFile:     c.String("rootwallet"),
				address:            c.String("address"),
				chamHashParamsFile: c.String("chamHashParams"),
			}

			return addAccount(args, logger)
		},
		Flags: []cli.Flag{
			headerFlag,
			feeFlag,
			rootkeyFlag,
			cli.StringFlag{
				Name:  "address",
				Usage: "the account's address",
			},
			cli.StringFlag{
				Name:  "chamHashParams",
				Usage: "save new chameleon hash parameters to `FILE`",
			},
		},
	}
}

func addAccount(args *addAccountArgs, logger *log.Logger) error {
	err := args.ValidateInput()
	if err != nil {
		return err
	}

	rootPrivKey, err := crypto.ExtractECDSAKeyFromFile(args.rootWalletFile)
	if err != nil {
		return err
	}

	chamHashParams, err := crypto.GetOrCreateChParamsFromFile(args.chamHashParamsFile)
	if err != nil {
		return err
	}

	chamHashCheckString := crypto.NewChCheckString(chamHashParams)

	var newAddress [64]byte
	newPubInt, _ := new(big.Int).SetString(args.address, 16)
	copy(newAddress[:], newPubInt.Bytes())

	tx, _, err := protocol.ConstrAccTx(
		byte(args.header),
		uint64(args.fee),
		newAddress,
		rootPrivKey,
		nil,
		nil,
		chamHashParams,
		chamHashCheckString,
		[]byte{},
	)
	if err != nil {
		return err
	}

	return sendAccountTx(tx, chamHashParams, logger)
}

func (args addAccountArgs) ValidateInput() error {
	if args.fee <= 0 {
		return errors.New("invalid argument: Fee must be > 0")
	}

	if len(args.rootWalletFile) == 0 {
		return errors.New("argument missing: rootwallet")
	}

	if len(args.address) == 0 {
		return errors.New("argument missing: address")
	}

	if len(args.address) != 128 {
		return errors.New("invalid argument length: address")
	}

	if len(args.chamHashParamsFile) == 0 {
		return errors.New("argument missing: chamHashParams")
	}

	return nil
}
