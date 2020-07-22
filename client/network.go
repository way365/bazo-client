package client

import (
	"errors"
	"github.com/julwil/bazo-client/args"
	"github.com/julwil/bazo-client/network"
	"github.com/julwil/bazo-client/util"
	"github.com/julwil/bazo-miner/crypto"
	"github.com/julwil/bazo-miner/p2p"
	"github.com/julwil/bazo-miner/protocol"
	"log"
)

func ConfigureNetwork(args *args.NetworkArgs, logger *log.Logger) error {
	err := args.ValidateInput()
	if err != nil {
		return err
	}

	privKey, err := crypto.ExtractECDSAKeyFromFile(args.TootWallet)
	if err != nil {
		return err
	}

	tx, err := protocol.ConstrConfigTx(
		byte(args.Header),
		uint8(args.OptionId),
		uint64(args.Payload),
		uint64(args.Fee),
		uint8(args.TxCount),
		privKey)

	if err != nil {
		return err
	}

	if tx == nil {
		return errors.New("transaction encoding failed")
	}

	if err := network.SendTx(util.Config.BootstrapIpport, tx, p2p.CONFIGTX_BRDCST); err != nil {
		logger.Printf("%v\n", err)
		return err
	} else {
		logger.Printf("Transaction successfully sent to network:\nTxHash: %x%v", tx.Hash(), tx)
	}

	return nil
}
