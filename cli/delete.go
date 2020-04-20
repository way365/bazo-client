package cli

import (
	"errors"
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
	txcount            int
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
				txcount:            c.Int("txcount"),
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
			cli.IntFlag{
				Name:  "txcount",
				Usage: "the sender's current transaction counter",
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

	// Finally, we create the delete-tx.
	tx, err := protocol.ConstrDeleteTx(
		byte(args.header),
		uint64(args.fee),
		txToDeleteHash,
		protocol.SerializeHashContent(issuerAddress),
		issuerPrivateKey,
		uint32(args.txcount),
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
