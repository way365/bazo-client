package services

import (
	"crypto/ecdsa"
	"crypto/rand"
	"github.com/way365/bazo-client/network"
	"github.com/way365/bazo-client/util"
	"github.com/way365/bazo-miner/p2p"
	"github.com/way365/bazo-miner/protocol"
	"log"
)

var (
	logger *log.Logger
)

func InitLogging() {
	logger = util.InitLogger()
}

func put(slice []*FundsTxJson, tx *FundsTxJson) {
	for i := 0; i < 9; i++ {
		slice[i] = slice[i+1]
	}

	slice[9] = tx
}

func SignTx(txHash [32]byte, tx protocol.Transaction, privKey *ecdsa.PrivateKey) error {
	var signature [64]byte
	r, s, err := ecdsa.Sign(rand.Reader, privKey, txHash[:])
	if err != nil {
		return err
	}

	copy(signature[:32], r.Bytes())
	copy(signature[32:], s.Bytes())
	tx.SetSignature(signature)

	return nil
}

func SubmitTx(txHash [32]byte, tx protocol.Transaction) error {
	var typeId uint8

	switch tx.(type) {
	case *protocol.AccTx:
		typeId = p2p.ACCTX_BRDCST
	case *protocol.FundsTx:
		typeId = p2p.FUNDSTX_BRDCST
	case *protocol.UpdateTx:
		typeId = p2p.UPDATETX_BRDCST
	case *protocol.StakeTx:
		typeId = p2p.STAKETX_BRDCST
	case *protocol.ConfigTx:
		typeId = p2p.CONFIGTX_BRDCST
	case *protocol.AggTx:
		typeId = p2p.AGGTX_BRDCST
	default:
		typeId = p2p.NOT_FOUND
	}

	if err := network.SendTx(util.Config.BootstrapIpport, tx, typeId); err != nil {
		logger.Printf("%v\n", err)
		return err
	}

	logger.Printf("Transaction successfully sent to network:\nTxHash: %x%v", txHash, tx.String())

	return nil
}
