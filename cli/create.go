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
	chParamsFile   string
	data           string
}

func getCreateAccountCommand(logger *log.Logger) cli.Command {
	return cli.Command{
		Name:  "create",
		Usage: "create a new account and add it to the network",
		Action: func(c *cli.Context) error {
			args := &createAccountArgs{
				header:         c.Int("Header"),
				fee:            c.Uint64("Fee"),
				rootWalletFile: c.String("rootwallet"),
				walletFile:     c.String("wallet"),
				chParamsFile:   c.String("chparams"),
				data:           c.String("Data"),
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
			cli.StringFlag{
				Name:  "chparams",
				Usage: "save new chameleon hash parameters to `FILE`",
			},
			cli.StringFlag{
				Name:  "Data",
				Usage: "Data field to add a message to the tx",
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

	chParams, err := crypto.GetOrCreateChParamsFromFile(args.chParamsFile)
	if err != nil {
		return err
	}

	// IMPORTANT: We need to sanitize the secret trapdoor key before we send it to the network.
	chParams.TK = []byte{}

	chCheckString := crypto.NewChCheckString(chParams)

	tx, newPrivKey, err := protocol.ConstrAccTx(
		byte(args.header),
		uint64(args.fee),
		[64]byte{},
		rootPrivKey,
		nil,
		nil,
		chParams,
		chCheckString,
		[]byte(args.data),
	)
	if err != nil {
		return err
	}

	//accAddress := protocol.SerializeHashContent(tx.PubKey)
	//crypto.ChamHashParamsMap[accAddress] = chParams

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

	return sendAccountTx(tx, chParams, logger)
}

func (args createAccountArgs) ValidateInput() error {
	if args.fee <= 0 {
		return errors.New("invalid argument: Fee must be > 0")
	}

	if len(args.rootWalletFile) == 0 {
		return errors.New("argument missing: rootWalletFile")
	}

	if len(args.walletFile) == 0 {
		return errors.New("argument missing: walletFile")
	}

	if len(args.chParamsFile) == 0 {
		return errors.New("argument missing: chamHashParams")
	}

	return nil
}
