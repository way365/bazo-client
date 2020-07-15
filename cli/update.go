package cli

import (
	"errors"
	"fmt"
	"github.com/julwil/bazo-client/cstorage"
	"github.com/julwil/bazo-client/network"
	"github.com/julwil/bazo-client/util"
	"github.com/julwil/bazo-miner/crypto"
	"github.com/julwil/bazo-miner/p2p"
	"github.com/julwil/bazo-miner/protocol"
	"github.com/urfave/cli"
	"log"
	"math/big"
)

type updateTxArgs struct {
	header             int
	fee                uint64
	txToUpdate         string
	txIssuerWalletFile string
	chParamsFile       string
	updateData         string // Data to be updated on the txToUpdate.
	data               string // Data on this tx.
}

func GetUpdateTxCommand(logger *log.Logger) cli.Command {
	return cli.Command{
		Name:  "update",
		Usage: "update the Data field of a specific transaction",
		Action: func(c *cli.Context) error {
			args := &updateTxArgs{
				header:             c.Int("header"),
				fee:                c.Uint64("fee"),
				txToUpdate:         c.String("tx-hash"),
				txIssuerWalletFile: c.String("tx-issuer"),
				chParamsFile:       c.String("chparams"),
				updateData:         c.String("update-Data"),
				data:               c.String("data"),
			}

			return updateTx(args, logger)
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

func updateTx(args *updateTxArgs, logger *log.Logger) error {
	err := args.ValidateInput()
	if err != nil {
		return err
	}

	// First, we read private key from creator from the wallet file
	issuerPrivateKey, err := crypto.ExtractECDSAKeyFromFile(args.txIssuerWalletFile)
	if err != nil {
		return err
	}

	// Then, we retrieve the associated address from that private key
	issuerAddress := crypto.GetAddressFromPubKey(&issuerPrivateKey.PublicKey)

	// Then, we parse the hash of the tx that shall be updated.
	var txToUpdateHash [32]byte
	if len(args.txToUpdate) == 64 {
		newPubInt, _ := new(big.Int).SetString(args.txToUpdate, 16)
		copy(txToUpdateHash[:], newPubInt.Bytes())
	}

	chParams, err := crypto.GetOrCreateChParamsFromFile(args.chParamsFile)
	chCheckString := crypto.NewChCheckString(chParams)
	if err != nil {
		return errors.New("no chameleon hash parameter files found with given parameters")
	}

	newData := []byte(args.updateData)
	// We create a new check string for TxToDelete to create a hash collision using chameleon hashing.
	newChCheckString := generateCollisionCheckString(txToUpdateHash, chParams, newData)

	// Finally, we create the update-tx.
	tx, err := protocol.ConstrUpdateTx(
		byte(args.header),
		uint64(args.fee),
		txToUpdateHash,
		newChCheckString,
		newData,
		protocol.SerializeHashContent(issuerAddress),
		issuerPrivateKey,
		chCheckString,
		chParams,
		[]byte(args.data),
	)

	if err != nil {
		logger.Printf("%v\n", err)
		return err
	}

	// Broadcast to the network
	if err := network.SendTx(util.Config.BootstrapIpport, tx, p2p.UPDATETX_BRDCST); err != nil {
		logger.Printf("%v\n", err)
		return err
	}

	txHash := tx.ChameleonHash(chParams)

	logger.Printf("Transaction successfully sent to network:\nTxHash: %x%v", txHash, tx)
	cstorage.WriteTransaction(txHash, tx)

	return nil
}

func (args updateTxArgs) ValidateInput() error {
	if len(args.txToUpdate) == 0 {
		return errors.New("argument missing: txHash")
	}

	if len(args.txIssuerWalletFile) == 0 {
		return errors.New("argument missing: txIssuer")
	}

	return nil
}

func generateCollisionCheckString(
	txToUpdateHash [32]byte,
	chParams *crypto.ChameleonHashParameters,
	newData []byte,
) (newChCheckString *crypto.ChameleonHashCheckString) {
	// First we need to query the Tx to update.
	var txToUpdate protocol.Transaction
	txToUpdate = cstorage.ReadTransaction(txToUpdateHash)
	if txToUpdate == nil {
		fmt.Printf("TX not found: %x", txToUpdateHash)

		return
	}

	fmt.Printf("TX to update %s", txToUpdate.String())

	// Then we have to save the old check string and the SHA3 hash before we mutate the tx.
	oldChCheckString := txToUpdate.GetChCheckString()
	oldSHA3 := txToUpdate.SHA3()
	oldHashInput := oldSHA3[:]

	// Now it's time to mutate the tx Data.
	txToUpdate.SetData(newData)

	// Then we compute the new SHA3 hash. This hash incorporates the changes in the Data field.
	// With the new hash input we compute a hash collision and get the new check string.
	newSHA3 := txToUpdate.SHA3()
	newHashInput := newSHA3[:]
	newChCheckString = crypto.GenerateChCollision(chParams, oldChCheckString, &oldHashInput, &newHashInput)

	// We update the txToUpdate record in our local db.
	txToUpdate.SetChCheckString(newChCheckString)
	cstorage.WriteTransaction(txToUpdateHash, txToUpdate)

	return newChCheckString
}
