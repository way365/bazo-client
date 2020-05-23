package cli

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/julwil/bazo-miner/crypto"
	"github.com/julwil/bazo-miner/protocol"
	"github.com/urfave/cli"
	"log"
	"os"
)

type createAccountArgs struct {
	header         int
	fee            uint64
	rootWalletFile string
	walletFile     string
}

func getCreateAccountCommand(logger *log.Logger) cli.Command {
	return cli.Command{
		Name:  "create",
		Usage: "create a new account and add it to the network",
		Action: func(c *cli.Context) error {
			args := &createAccountArgs{
				header:         c.Int("header"),
				fee:            c.Uint64("fee"),
				rootWalletFile: c.String("rootwallet"),
				walletFile:     c.String("wallet"),
			}

			return createAccount(args, logger)
		},
		Flags: []cli.Flag{
			headerFlag,
			feeFlag,
			rootkeyFlag,
			cli.StringFlag{
				Name:  "wallet",
				Usage: "save new account's public private key to `FILE`",
			},
		},
	}
}

func createAccount(args *createAccountArgs, logger *log.Logger) error {
	err := args.ValidateInput()
	if err != nil {
		return err
	}

	rootPrivKey, err := crypto.ExtractECDSAKeyFromFile(args.rootWalletFile)
	if err != nil {
		return err
	}

	// Private Key
	var newPrivKey *ecdsa.PrivateKey

	tx, newPrivKey, err := protocol.ConstrAccTx(
		byte(args.header),
		uint64(args.fee),
		[64]byte{},
		rootPrivKey,
		nil,
		nil,
		nil,
	)
	if err != nil {
		return err
	}

	//Write the private key to the given textfile
	file, err := os.Create(args.walletFile)
	if err != nil {
		return err
	}

	_, err = file.WriteString(string(newPrivKey.X.Text(16)) + "\n")
	_, err = file.WriteString(string(newPrivKey.Y.Text(16)) + "\n")
	_, err = file.WriteString(string(newPrivKey.D.Text(16)) + "\n")

	if err != nil {
		return errors.New(fmt.Sprintf("failed to write key to file %v", args.walletFile))
	}

	return sendAccountTx(tx, logger)
}

func (args createAccountArgs) ValidateInput() error {
	if args.fee <= 0 {
		return errors.New("invalid argument: fee must be > 0")
	}

	if len(args.rootWalletFile) == 0 {
		return errors.New("argument missing: rootWalletFile")
	}

	if len(args.walletFile) == 0 {
		return errors.New("argument missing: walletFile")
	}

	return nil
}
