package args

import "errors"

type NetworkArgs struct {
	Header     int
	Fee        uint64
	TxCount    int
	TootWallet string
	OptionId   uint8
	Payload    uint64
}

type ConfigOption struct {
	Id    uint8
	Name  string
	Usage string
}

func (args NetworkArgs) ValidateInput() error {
	if args.Fee <= 0 {
		return errors.New("invalid argument: Fee must be > 0")
	}

	if args.TxCount < 0 {
		return errors.New("invalid argument: txcnt must be >= 0")
	}

	if len(args.TootWallet) == 0 {
		return errors.New("argument missing: rootwallet")
	}

	return nil
}
