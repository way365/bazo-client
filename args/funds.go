package args

import "errors"

type FundsArgs struct {
	Header         int    `json:"header"`
	FromWalletFile string `json:"from"`
	ToWalletFile   string `json:"to"`
	ToAddress      string
	MultisigFile   string `json:"multisig"`
	ChParamsFile   string `json:"chparams"`
	Amount         uint64 `json:"amount"`
	Fee            uint64 `json:"fee"`
	TxCount        int    `json:"txcount"`
	Data           string `json:"data"`
}

func (args FundsArgs) ValidateInput() error {
	if len(args.FromWalletFile) == 0 {
		return errors.New("argument missing: from")
	}

	if args.TxCount < 0 {
		return errors.New("invalid argument: txcnt must be >= 0")
	}

	if len(args.ToWalletFile) == 0 && len(args.ToAddress) == 0 {
		return errors.New("argument missing: to or toAddess")
	}

	if len(args.ToWalletFile) == 0 && len(args.ToAddress) != 128 {
		return errors.New("invalid argument: ToAddress")
	}

	if args.Fee <= 0 {
		return errors.New("invalid argument: Fee must be > 0")
	}

	if args.Amount <= 0 {
		return errors.New("invalid argument: Amount must be > 0")
	}

	return nil
}
