package services

import (
	"encoding/hex"
	"github.com/way365/bazo-miner/protocol"
)

type FundsTxJson struct {
	Header byte   `json:"header"`
	Hash   string `json:"hash"`
	Amount uint64 `json:"amount"`
	Fee    uint64 `json:"fee"`
	TxCnt  uint32 `json:"txCnt"`
	From   string `json:"from"`
	To     string `json:"to"`
	Sig1   string `json:"sig1"`
	Sig2   string `json:"sig2"`
	Status string `json:"status"`
}

func ConvertFundsTx(fundsTx *protocol.FundsTx, status string) (fundsTxJson *FundsTxJson) {
	txHash := fundsTx.Hash()
	return &FundsTxJson{
		fundsTx.Header,
		hex.EncodeToString(txHash[:]),
		fundsTx.Amount,
		fundsTx.Fee,
		fundsTx.TxCnt,
		hex.EncodeToString(fundsTx.From[:]),
		hex.EncodeToString(fundsTx.To[:]),
		hex.EncodeToString(fundsTx.Sig1[:]),
		hex.EncodeToString(fundsTx.Sig2[:]),
		status,
	}
}
