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
	chamHashParamsFile string
	updateData         string
}

func GetUpdateTxCommand(logger *log.Logger) cli.Command {
	return cli.Command{
		Name:  "update",
		Usage: "update the data field of a specific transaction",
		Action: func(c *cli.Context) error {
			args := &updateTxArgs{
				header:             c.Int("header"),
				fee:                c.Uint64("fee"),
				txToUpdate:         c.String("tx-hash"),
				txIssuerWalletFile: c.String("tx-issuer"),
				chamHashParamsFile: c.String("cham-hash-params"),
				updateData:         c.String("update-data"),
			}

			return updateTx(args, logger)
		},
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "header",
				Usage: "header flag",
				Value: 0,
			},
			cli.Uint64Flag{
				Name:  "fee",
				Usage: "specify the fee",
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
				Name:  "cham-hash-params",
				Usage: "load the chameleon hash parameters from `FILE`",
			},
			cli.StringFlag{
				Name:  "update-data",
				Usage: "specify the new data that shall be updated on the tx",
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

	chamHashParams, err := crypto.GetOrCreateChamHashParamsFromFile(args.chamHashParamsFile)
	if err != nil {
		return errors.New("no chameleon hash parameter files found with given parameters")
	}

	newData := []byte(args.updateData)
	// We create a new check string for TxToDelete to create a hash collision using chameleon hashing.
	newChamHashCheckString := generateCollisionCheckString(txToUpdateHash, chamHashParams, newData)

	// Finally, we create the update-tx.
	tx, err := protocol.ConstrUpdateTx(
		byte(args.header),
		uint64(args.fee),
		txToUpdateHash,
		newChamHashCheckString,
		newData,
		protocol.SerializeHashContent(issuerAddress),
		issuerPrivateKey,
	)

	if err != nil {
		logger.Printf("%v\n", err)
		return err
	}

	// Broadcast to the network
	if err := network.SendTx(util.Config.BootstrapIpport, tx, p2p.UPDATETX_BRDCST); err != nil {
		logger.Printf("%v\n", err)
		return err
	} else {
		logger.Printf("Transaction successfully sent to network:\nTxHash: %x%v", tx.Hash(), tx)
	}

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
	txToDeleteHash [32]byte,
	parameters *crypto.ChameleonHashParameters,
	newData []byte,
) (newCheckString *crypto.ChameleonHashCheckString) {
	// First we need to query the Tx to update.
	var txToUpdate protocol.Transaction
	txToUpdate = cstorage.ReadTransaction(txToDeleteHash)

	fmt.Printf("TX to update %s", txToUpdate.String())

	// Then we have to save the old check string and the SHA3 hash before we mutate the tx.
	oldCheckString := txToUpdate.GetChamHashCheckString()
	oldSHA3 := txToUpdate.SHA3()
	oldHashInput := oldSHA3[:]

	// Now it's time to mutate the tx data.
	txToUpdate.SetData(newData)

	// Then we compute the new SHA3 hash. This hash incorporates the changes in the data field.
	// With the new hash input we compute a hash collision and get the new check string.
	newSHA3 := txToUpdate.SHA3()
	newHashInput := newSHA3[:]
	newCheckString = crypto.GenerateChamHashCollision(parameters, oldCheckString, &oldHashInput, &newHashInput)

	// We update the tx record in our local db.
	txToUpdate.SetChamHashCheckString(newCheckString)

	//fmt.Printf("\nAFTER update (%x): %s", txToUpdate.HashWithChamHashParams(parameters), txToUpdate.String())

	cstorage.WriteTransaction(txToDeleteHash, txToUpdate)

	return newCheckString
}
