package args

import (
	"errors"
	"github.com/urfave/cli"
)

type StakingArgs struct {
	Header       int
	Fee          uint64
	Wallet       string
	Commitment   string
	StakingValue bool
}

func ParseStakingArgs(c *cli.Context) *StakingArgs {
	return &StakingArgs{
		Header:     c.Int("header"),
		Fee:        c.Uint64("fee"),
		Wallet:     c.String("wallet"),
		Commitment: c.String("commitment"),
	}
}

func (args StakingArgs) ValidateInput() error {
	if args.Fee <= 0 {
		return errors.New("invalid argument: Fee must be > 0")
	}

	if len(args.Wallet) == 0 {
		return errors.New("argument missing: wallet")
	}

	if args.StakingValue && len(args.Commitment) == 0 {
		return errors.New("argument missing: Commitment")
	}

	return nil
}
