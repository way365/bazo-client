package args

import "errors"

type UpdateTxArgs struct {
	Header         int
	Fee            uint64
	TxToUpdate     string
	TxIssuerWallet string
	ChParams       string
	UpdateData     string // Data to be updated on the TxToUpdate.
	Data           string // Data on this tx.
}

func (args UpdateTxArgs) ValidateInput() error {
	if len(args.TxToUpdate) == 0 {
		return errors.New("argument missing: txHash")
	}

	if len(args.TxIssuerWallet) == 0 {
		return errors.New("argument missing: txIssuer")
	}

	return nil
}
