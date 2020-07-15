package client

import (
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/julwil/bazo-client/args"
	"github.com/julwil/bazo-client/cstorage"
	"github.com/julwil/bazo-client/network"
	"github.com/julwil/bazo-client/util"
	"github.com/julwil/bazo-miner/crypto"
	"github.com/julwil/bazo-miner/p2p"
	"github.com/julwil/bazo-miner/protocol"
	"log"
)

func PrepareFundsTx(args *args.FundsArgs, logger *log.Logger) (txHash [32]byte, err error) {
	fromPrivKey, err := crypto.ExtractECDSAKeyFromFile(args.FromWalletFile)
	if err != nil {
		return [32]byte{}, err
	}

	var toPubKey *ecdsa.PublicKey
	if len(args.ToWalletFile) == 0 {
		if len(args.ToAddress) == 0 {
			return [32]byte{}, errors.New(fmt.Sprintln("No recipient specified"))
		} else {
			if len(args.ToAddress) != 128 {
				return [32]byte{}, errors.New(fmt.Sprintln("Invalid recipient address"))
			}

			runes := []rune(args.ToAddress)
			pub1 := string(runes[:64])
			pub2 := string(runes[64:])

			toPubKey, err = crypto.GetPubKeyFromString(pub1, pub2)
			if err != nil {
				return [32]byte{}, err
			}
		}
	} else {
		toPubKey, err = crypto.ExtractECDSAPublicKeyFromFile(args.ToWalletFile)
		if err != nil {
			return [32]byte{}, err
		}
	}

	var multisigPrivKey *ecdsa.PrivateKey
	if len(args.MultisigFile) > 0 {
		multisigPrivKey, err = crypto.ExtractECDSAKeyFromFile(args.MultisigFile)
		if err != nil {
			return [32]byte{}, err
		}
	} else {
		multisigPrivKey = fromPrivKey
	}

	fromAddress := crypto.GetAddressFromPubKey(&fromPrivKey.PublicKey)
	toAddress := crypto.GetAddressFromPubKey(toPubKey)

	chParams, err := crypto.GetOrCreateChParamsFromFile(args.ChParamsFile)
	chCheckString := crypto.NewChCheckString(chParams)

	tx, err := protocol.ConstrFundsTx(
		byte(args.Header),
		uint64(args.Amount),
		uint64(args.Fee),
		uint32(args.TxCount),
		protocol.SerializeHashContent(fromAddress),
		protocol.SerializeHashContent(toAddress),
		fromPrivKey,
		multisigPrivKey,
		chCheckString,
		chParams,
		[]byte(args.Data),
	)

	if err != nil {
		logger.Printf("%v\n", err)
		return [32]byte{}, err
	}

	txHash = tx.ChameleonHash(chParams)
	cstorage.WriteTransaction(txHash, tx)

	return txHash, nil
}

func SignFundsTx(
	tx *protocol.FundsTx,
	privKey1 *ecdsa.PrivateKey,
	privKey2 *ecdsa.PrivateKey,
	chParams *crypto.ChameleonHashParameters,
) error {
	txHash := tx.ChameleonHash(chParams)

	r, s, err := ecdsa.Sign(rand.Reader, privKey1, txHash[:])
	if err != nil {
		return err
	}

	copy(tx.Sig1[32-len(r.Bytes()):32], r.Bytes())
	copy(tx.Sig1[64-len(s.Bytes()):], s.Bytes())

	if privKey2 != nil {
		r, s, err := ecdsa.Sign(rand.Reader, privKey2, txHash[:])
		if err != nil {
			return err
		}

		copy(tx.Sig2[32-len(r.Bytes()):32], r.Bytes())
		copy(tx.Sig2[64-len(s.Bytes()):], s.Bytes())
	}

	return nil
}

func SubmitFundsTx(tx *protocol.FundsTx) error {
	err := network.SendTx(util.Config.BootstrapIpport, tx, p2p.FUNDSTX_BRDCST)

	return err
}

func CreateSignSubmitFundsTx(args *args.FundsArgs, logger *log.Logger) (txHash [32]byte, err error) {
	fromPrivKey, err := crypto.ExtractECDSAKeyFromFile(args.FromWalletFile)
	if err != nil {
		return [32]byte{}, err
	}

	var toPubKey *ecdsa.PublicKey
	if len(args.ToWalletFile) == 0 {
		if len(args.ToAddress) == 0 {
			return [32]byte{}, errors.New(fmt.Sprintln("No recipient specified"))
		} else {
			if len(args.ToAddress) != 128 {
				return [32]byte{}, errors.New(fmt.Sprintln("Invalid recipient address"))
			}

			runes := []rune(args.ToAddress)
			pub1 := string(runes[:64])
			pub2 := string(runes[64:])

			toPubKey, err = crypto.GetPubKeyFromString(pub1, pub2)
			if err != nil {
				return [32]byte{}, err
			}
		}
	} else {
		toPubKey, err = crypto.ExtractECDSAPublicKeyFromFile(args.ToWalletFile)
		if err != nil {
			return [32]byte{}, err
		}
	}

	var multisigPrivKey *ecdsa.PrivateKey
	if len(args.MultisigFile) > 0 {
		multisigPrivKey, err = crypto.ExtractECDSAKeyFromFile(args.MultisigFile)
		if err != nil {
			return [32]byte{}, err
		}
	} else {
		multisigPrivKey = fromPrivKey
	}

	fromAddress := crypto.GetAddressFromPubKey(&fromPrivKey.PublicKey)
	toAddress := crypto.GetAddressFromPubKey(toPubKey)

	chParams, err := crypto.GetOrCreateChParamsFromFile(args.ChParamsFile)
	chCheckString := crypto.NewChCheckString(chParams)

	tx, err := protocol.ConstrFundsTx(
		byte(args.Header),
		uint64(args.Amount),
		uint64(args.Fee),
		uint32(args.TxCount),
		protocol.SerializeHashContent(fromAddress),
		protocol.SerializeHashContent(toAddress),
		fromPrivKey,
		multisigPrivKey,
		chCheckString,
		chParams,
		[]byte(args.Data),
	)

	if err != nil {
		logger.Printf("%v\n", err)
		return [32]byte{}, err
	}

	if err := SignFundsTx(tx, fromPrivKey, multisigPrivKey, chParams); err != nil {
		logger.Printf("%v\n", err)
		return [32]byte{}, err
	}

	if err := SubmitFundsTx(tx); err != nil {
		logger.Printf("%v\n", err)
		return [32]byte{}, err
	}

	txHash = tx.ChameleonHash(chParams)

	logger.Printf("Transaction successfully sent to network:\nTxHash: %x%v", txHash, tx)
	cstorage.WriteTransaction(txHash, tx)

	return txHash, nil
}
