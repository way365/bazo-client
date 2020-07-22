package client

import (
	"crypto/rsa"
	"errors"
	"github.com/julwil/bazo-client/args"
	"github.com/julwil/bazo-client/network"
	"github.com/julwil/bazo-client/util"
	"github.com/julwil/bazo-miner/crypto"
	"github.com/julwil/bazo-miner/p2p"
	"github.com/julwil/bazo-miner/protocol"
	"log"
)

func ToggleStaking(args *args.StakingArgs, logger *log.Logger) error {
	err := args.ValidateInput()
	if err != nil {
		return err
	}

	privKey, err := crypto.ExtractECDSAKeyFromFile(args.Wallet)
	if err != nil {
		return err
	}

	accountPubKey := crypto.GetAddressFromPubKey(&privKey.PublicKey)

	commPubKey := &rsa.PublicKey{}
	if args.StakingValue {
		commPrivKey, err := crypto.ExtractRSAKeyFromFile(args.Commitment)
		if err != nil {
			return err
		}
		commPubKey = &commPrivKey.PublicKey
	}

	tx, err := protocol.ConstrStakeTx(
		byte(args.Header),
		uint64(args.Fee),
		args.StakingValue,
		protocol.SerializeHashContent(accountPubKey),
		privKey,
		commPubKey,
	)

	if err != nil {
		return err
	}

	if tx == nil {
		return errors.New("transaction encoding failed")
	}

	if err := network.SendTx(util.Config.BootstrapIpport, tx, p2p.STAKETX_BRDCST); err != nil {
		logger.Printf("%v\n", err)
		return err
	} else {
		logger.Printf("Transaction successfully sent to network:\nTxHash: %x%v", tx.Hash(), tx)
	}

	return nil
}
