package args

import "errors"

type UpdateTxArgs struct {
	Header     int    `json:"header"`
	Fee        uint64 `json:"fee"`
	TxToUpdate string `json:"tx_to_update"`
	TxIssuer   string `json:"tx_issuer"`
	ChParams   string `json:"ch_params"`
	UpdateData string `json:"update_data"`
	Data       string `json:"data"`
}

func (args UpdateTxArgs) ValidateInput() error {
	if len(args.TxToUpdate) == 0 {
		return errors.New("argument missing: txHash")
	}

	if len(args.TxIssuer) == 0 {
		return errors.New("argument missing: txIssuer")
	}

	return nil
}
