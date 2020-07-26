package args

import (
	"errors"
)

type FundsArgs struct {
	Header      int    `json:"header"`
	From        string `json:"from"`
	To          string `json:"to"`
	MultiSigKey string `json:"multi_sig"`
	ChParams    string `json:"ch_params"`
	Amount      uint64 `json:"amount"`
	Fee         uint64 `json:"fee"`
	TxCount     int    `json:"tx_count"`
	Data        string `json:"data"`
}

func (args FundsArgs) ValidateInput() error {
	if len(args.From) == 0 {
		return errors.New("argument missing: from")
	}

	if args.TxCount < 0 {
		return errors.New("invalid argument: txcnt must be >= 0")
	}

	if len(args.To) == 0 {
		return errors.New("argument missing: to")
	}

	if args.Fee <= 0 {
		return errors.New("invalid argument: Fee must be > 0")
	}

	if args.Amount <= 0 {
		return errors.New("invalid argument: Amount must be > 0")
	}

	return nil
}
