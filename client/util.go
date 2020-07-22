package client

import (
	"github.com/julwil/bazo-client/network"
	"github.com/julwil/bazo-client/util"
	"github.com/julwil/bazo-miner/p2p"
	"github.com/julwil/bazo-miner/protocol"
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
