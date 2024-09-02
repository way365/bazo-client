package services

import (
	"crypto/ecdsa"
	"crypto/rand"
	"github.com/way365/bazo-client/args"
	"github.com/way365/bazo-client/cstorage"
	"github.com/way365/bazo-miner/crypto"
	"github.com/way365/bazo-miner/protocol"
	"log"
)

func PrepareSignSubmitFundsTx(arguments *args.FundsArgs, logger *log.Logger) (txHash [32]byte, err error) {
	err = arguments.ValidateInput()
	if err != nil {
		return [32]byte{}, err
	}

	txHash, tx, err := PrepareFundsTx(arguments, logger)

	fromPrivKey, err := args.ResolvePrivateKey(arguments.From)
	if err != nil {
		return [32]byte{}, err
	}

	multiSigPrivKey, err := args.ResolvePrivateKey(arguments.MultiSigKey)
	if err != nil {
		return [32]byte{}, err
	}

	if err := SignFundsTx(txHash, tx, fromPrivKey, multiSigPrivKey); err != nil {
		logger.Printf("%v\n", err)
		return [32]byte{}, err
	}

	if err := SubmitTx(txHash, tx); err != nil {
		logger.Printf("%v\n", err)
		return [32]byte{}, err
	}

	cstorage.WriteTransaction(txHash, tx)

	return txHash, nil
}

func PrepareFundsTx(arguments *args.FundsArgs, logger *log.Logger) (txHash [32]byte, tx *protocol.FundsTx, err error) {
	err = arguments.ValidateInput()
	if err != nil {
		return [32]byte{}, tx, err
	}

	fromPubKey, err := args.ResolvePublicKey(arguments.From)
	if err != nil {
		return [32]byte{}, tx, err
	}

	toPubKey, err := args.ResolvePublicKey(arguments.To)
	if err != nil {
		return [32]byte{}, tx, err
	}

	fromAddress := crypto.GetAddressFromPubKey(fromPubKey)
	toAddress := crypto.GetAddressFromPubKey(toPubKey)

	parameters, err := args.ResolveParameters(arguments.Parameters)
	if err != nil {
		return [32]byte{}, tx, err
	}

	checkString := crypto.NewCheckString(parameters)

	tx, err = protocol.ConstrFundsTx(
		byte(arguments.Header),
		uint64(arguments.Amount),
		uint64(arguments.Fee),
		uint32(arguments.TxCount),
		protocol.SerializeHashContent(fromAddress),
		protocol.SerializeHashContent(toAddress),
		checkString,
		[]byte(arguments.Data),
	)

	if err != nil {
		logger.Printf("%v\n", err)
		return [32]byte{}, tx, err
	}

	txHash = tx.ChameleonHash(parameters)
	cstorage.WriteTransaction(txHash, tx)

	return txHash, tx, err
}

func SignFundsTx(txHash [32]byte, tx *protocol.FundsTx, privKey *ecdsa.PrivateKey, multiSigKey *ecdsa.PrivateKey) error {
	r, s, err := ecdsa.Sign(rand.Reader, privKey, txHash[:])
	if err != nil {
		return err
	}

	copy(tx.Sig1[32-len(r.Bytes()):32], r.Bytes())
	copy(tx.Sig1[64-len(s.Bytes()):], s.Bytes())

	if multiSigKey != nil {
		r, s, err := ecdsa.Sign(rand.Reader, multiSigKey, txHash[:])
		if err != nil {
			return err
		}

		copy(tx.Sig2[32-len(r.Bytes()):32], r.Bytes())
		copy(tx.Sig2[64-len(s.Bytes()):], s.Bytes())
	}

	return nil
}
