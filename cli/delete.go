package cli

import (
	"errors"
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

type deleteTxArgs struct {
	header             int
	fee                uint64
	txToDeleteHash     string
	txIssuerWalletFile string
	chamHashParamsFile string
	newData            string
}

func GetDeleteTxCommand(logger *log.Logger) cli.Command {
	return cli.Command{
		Name:  "delete-tx",
		Usage: "delete a specific transaction",
		Action: func(c *cli.Context) error {
			args := &deleteTxArgs{
				header:             c.Int("header"),
				fee:                c.Uint64("fee"),
				txToDeleteHash:     c.String("txHash"),
				txIssuerWalletFile: c.String("txIssuer"),
				chamHashParamsFile: c.String("chamHashParams"),
				newData:            c.String("newData"),
			}

			return deleteTx(args, logger)
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
				Name:  "txHash",
				Usage: "the 32-byte hash of the transaction to delete",
			},
			cli.StringFlag{
				Name:  "txIssuer",
				Usage: "load the tx issuer's public key from `FILE`",
			},
			cli.StringFlag{
				Name:  "chamHashParams",
				Usage: "load the chameleon hash parameters from `FILE`",
			},
			cli.StringFlag{
				Name:  "newData",
				Usage: "specify the new data that shall be updated on the tx",
			},
		},
	}
}

func deleteTx(args *deleteTxArgs, logger *log.Logger) error {
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

	// Then, we parse the hash of the tx that shall be deleted.
	var txToDeleteHash [32]byte
	if len(args.txToDeleteHash) == 64 {
		newPubInt, _ := new(big.Int).SetString(args.txToDeleteHash, 16)
		copy(txToDeleteHash[:], newPubInt.Bytes())
	}

	chamHashParams, err := crypto.GetOrCreateChamHashParamsFromFile(args.chamHashParamsFile)
	if err != nil {
		return errors.New("no chameleon hash parameter files found with given parameters")
	}

	newData := []byte(args.newData)
	// We create a new check string for TxToDelete to create a hash collision using chameleon hashing.
	newChamHashCheckString := generateCollisionCheckString(txToDeleteHash, chamHashParams, newData)

	// Finally, we create the delete-tx.
	tx, err := protocol.ConstrDeleteTx(
		byte(args.header),
		uint64(args.fee),
		txToDeleteHash,
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
	if err := network.SendTx(util.Config.BootstrapIpport, tx, p2p.DELTX_BRDCST); err != nil {
		logger.Printf("%v\n", err)
		return err
	} else {
		logger.Printf("Transaction successfully sent to network:\nTxHash: %x%v", tx.Hash(), tx)
	}

	return nil
}

func (args deleteTxArgs) ValidateInput() error {
	if len(args.txToDeleteHash) == 0 {
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
	// First we need to query the Tx to delete.
	var txToDelete protocol.Transaction
	txToDelete = cstorage.ReadTransaction(txToDeleteHash)

	//fmt.Printf("BEFORE update (%x): %s", txToDelete.HashWithChamHashParams(parameters), txToDelete.String())

	// Then we have to save the old check string and the SHA3 hash before we mutate the tx.
	oldCheckString := txToDelete.GetChamHashCheckString()
	oldSHA3 := txToDelete.SHA3()
	oldHashInput := oldSHA3[:]

	// Now it's time to mutate the tx data.
	txToDelete.SetData(newData)

	// Then we compute the new SHA3 hash. This hash incorporates the changes in the data field.
	// With the new hash input we compute a hash collision and get the new check string.
	newSHA3 := txToDelete.SHA3()
	newHashInput := newSHA3[:]
	newCheckString = crypto.GenerateChamHashCollision(parameters, oldCheckString, &oldHashInput, &newHashInput)

	// We update the tx record in our local db.
	txToDelete.SetChamHashCheckString(newCheckString)

	//fmt.Printf("\nAFTER update (%x): %s", txToDelete.HashWithChamHashParams(parameters), txToDelete.String())

	cstorage.WriteTransaction(txToDeleteHash, txToDelete)

	return newCheckString
}
