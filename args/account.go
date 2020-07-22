package args

import "errors"

type CreateAccountArgs struct {
	Header     int    `json:"header"`
	Fee        uint64 `json:"fee"`
	RootWallet string `json:"root_wallet"`
	Wallet     string `json:"wallet"`
	ChParams   string `json:"chparams"`
	Data       string `json:"data"`
}

type AddAccountArgs struct {
	Header     int    `json:"header"`
	Fee        uint64 `json:"fee"`
	RootWallet string `json:"root_wallet"`
	Address    string `json:"address"`
	ChParams   string `json:"chparams"`
}

type CheckAccountArgs struct {
	Address string `json:"address"`
	Wallet  string `json:"wallet"`
}

func (args CreateAccountArgs) ValidateInput() error {
	if args.Fee <= 0 {
		return errors.New("invalid argument: Fee must be > 0")
	}

	if len(args.RootWallet) == 0 {
		return errors.New("argument missing: rootwallet")
	}

	if len(args.Wallet) == 0 {
		return errors.New("argument missing: wallet")
	}

	if len(args.ChParams) == 0 {
		return errors.New("argument missing: chparams")
	}

	return nil
}

func (args AddAccountArgs) ValidateInput() error {
	if args.Fee <= 0 {
		return errors.New("invalid argument: Fee must be > 0")
	}

	if len(args.RootWallet) == 0 {
		return errors.New("argument missing: rootwallet")
	}

	if len(args.Address) == 0 {
		return errors.New("argument missing: Address")
	}

	if len(args.Address) != 128 {
		return errors.New("invalid argument length: Address")
	}

	if len(args.ChParams) == 0 {
		return errors.New("argument missing: chparams")
	}

	return nil
}

func (args CheckAccountArgs) ValidateInput() error {
	if len(args.Address) == 0 && len(args.Wallet) == 0 {
		return errors.New("argument missing: Address or wallet")
	}

	if len(args.Wallet) == 0 && len(args.Address) != 128 {
		return errors.New("invalid argument: Address")
	}

	return nil
}
