package client

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/julwil/bazo-client/args"
	"github.com/julwil/bazo-client/cstorage"
	"github.com/julwil/bazo-client/network"
	"github.com/julwil/bazo-miner/crypto"
	"github.com/julwil/bazo-miner/miner"
	"github.com/julwil/bazo-miner/protocol"
	"log"
	"math/big"
)

type Account struct {
	Address       [64]byte `json:"-"`
	AddressString string   `json:"address"`
	Balance       uint64   `json:"balance"`
	TxCnt         uint32   `json:"txCnt"`
	IsCreated     bool     `json:"isCreated"`
	IsRoot        bool     `json:"isRoot"`
	IsStaking     bool     `json:"isStaking"`
}

func PrepareSignSubmitCreateAccTx(arguments *args.CreateAccountArgs, logger *log.Logger) (txHash [32]byte, err error) {
	txHash, tx, err := PrepareCreateAccountTx(arguments, logger)
	if err != nil {
		return [32]byte{}, err
	}

	issuerPrivKey, err := args.ResolvePrivateKey(arguments.RootWallet)
	if err != nil {
		return [32]byte{}, err
	}

	if err := SignAccountTx(txHash, tx, issuerPrivKey); err != nil {
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

func PrepareCreateAccountTx(arguments *args.CreateAccountArgs, logger *log.Logger) (txHash [32]byte, tx *protocol.AccTx, err error) {
	err = arguments.ValidateInput()
	if err != nil {
		return [32]byte{}, tx, err
	}

	newPubKey, err := crypto.GetOrCreateECDSAPublicKeyFromFile(arguments.Wallet)
	if err != nil {
		return [32]byte{}, tx, err

	}

	newChParams, err := crypto.GetOrCreateChParamsFromFile(arguments.ChParams)
	if err != nil {
		return [32]byte{}, tx, err
	}

	// IMPORTANT: We need to sanitize the secret trapdoor key before we send it to the network.
	newChParams.TK = []byte{}

	chCheckString := crypto.NewChCheckString(newChParams)

	issuerPubKey, err := args.ResolvePublicKey(arguments.RootWallet)
	if err != nil {
		return [32]byte{}, tx, err
	}

	tx, err = protocol.ConstrAccTx(
		byte(arguments.Header),
		uint64(arguments.Fee),
		protocol.SerializeHashContent(crypto.GetAddressFromPubKey(issuerPubKey)),
		crypto.GetAddressFromPubKey(newPubKey),
		nil,
		nil,
		newChParams,
		chCheckString,
		[]byte(arguments.Data),
	)
	if err != nil {
		return [32]byte{}, tx, err
	}

	txHash = tx.ChameleonHash(newChParams)
	cstorage.WriteTransaction(txHash, tx)

	return txHash, tx, err
}

func SignAccountTx(txHash [32]byte, tx *protocol.AccTx, privKey *ecdsa.PrivateKey) error {
	var signature [64]byte
	r, s, err := ecdsa.Sign(rand.Reader, privKey, txHash[:])
	if err != nil {
		return err
	}

	copy(signature[:32], r.Bytes())
	copy(signature[32:], s.Bytes())
	tx.SetSignature(signature)

	return nil
}

//func SubmitAccountTx(txHash [32]byte, tx protocol.Transaction) error {
//
//	if err := network.SendTx(util.Config.BootstrapIpport, tx, p2p.ACCTX_BRDCST); err != nil {
//		logger.Printf("%v\n", err)
//		return err
//	} else {
//		logger.Printf("Transaction successfully sent to network:\nTxHash: %x%v", txHash, tx)
//	}
//
//	return nil
//}

func GetAccount(address [64]byte) (*Account, []*FundsTxJson, error) {
	//Initialize new account with empty address
	account := Account{address, hex.EncodeToString(address[:]), 0, 0, false, false, false}

	//Set default params
	activeParameters = miner.NewDefaultParameters()

	network.AccReq(false, protocol.SerializeHashContent(account.Address))
	if accI, _ := network.Fetch(network.AccChan); accI != nil {
		if acc := accI.(*protocol.Account); acc != nil {
			account.IsCreated = true
			account.IsStaking = acc.IsStaking

			//If Acc is Root in the bazo network state, we do not check for accTx, else we check
			network.AccReq(true, protocol.SerializeHashContent(account.Address))
			if rootAccI, _ := network.Fetch(network.AccChan); rootAccI != nil {
				if rootAcc := rootAccI.(*protocol.Account); rootAcc != nil {
					account.IsRoot = true
				}
			}
		}
	}

	if account.IsCreated == false {
		return nil, nil, errors.New(fmt.Sprintf("Account %x does not exist.\n", account.Address[:8]))
	}

	//if account.IsStaking == true {
	//	return nil, nil, errors.New(fmt.Sprintf("Account %x is a validator account. Validator's state cannot be calculated at the moment. We are sorry.\n", account.Address[:8]))
	//}

	var lastTenTx = make([]*FundsTxJson, 10)
	err := getState(&account, lastTenTx)
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("Could not calculate state of account %x: %v\n", account.Address[:8], err))
	}

	//No accTx exists for this account since it is the initial root account
	//Add the initial root's balance
	//if account.IsCreated == false && account.IsRoot == true {
	//	account.IsCreated = true
	//}

	return &account, lastTenTx, nil
}

func AddAccount(arguments *args.AddAccountArgs, logger *log.Logger) error {
	err := arguments.ValidateInput()
	if err != nil {
		return err
	}

	rootPrivKey, err := args.ResolvePrivateKey(arguments.RootWallet)
	if err != nil {
		return err
	}

	chParams, err := args.ResolveChParams(arguments.ChParams)
	if err != nil {
		return err
	}

	chCheckString := crypto.NewChCheckString(chParams)

	var addressBytes [64]byte
	copy(addressBytes[:], arguments.Address)

	tx, err := protocol.ConstrAccTx(
		byte(arguments.Header),
		uint64(arguments.Fee),
		protocol.SerializeHashContent(crypto.GetAddressFromPubKey(&rootPrivKey.PublicKey)),
		addressBytes,
		nil,
		nil,
		chParams,
		chCheckString,
		[]byte{},
	)
	if err != nil {
		return err
	}

	txHash := tx.ChameleonHash(chParams)

	return SubmitTx(txHash, tx)
}

func CheckAccount(args *args.CheckAccountArgs, logger *log.Logger) error {
	err := args.ValidateInput()
	if err != nil {
		return err
	}

	var address [64]byte
	if len(args.Address) == 128 {
		newPubInt, _ := new(big.Int).SetString(args.Address, 16)
		copy(address[:], newPubInt.Bytes())
	} else {
		privKey, err := crypto.ExtractECDSAKeyFromFile(args.Wallet)
		if err != nil {
			logger.Printf("%v\n", err)
			return err
		}

		address = crypto.GetAddressFromPubKey(&privKey.PublicKey)
	}

	logger.Printf("My Address: %x\n", address)

	loadBlockHeaders()
	acc, _, err := GetAccount(address)
	if err != nil {
		logger.Println(err)
		return err
	} else {
		logger.Printf(acc.String())
	}

	return nil
}

func (acc Account) String() string {
	addressHash := protocol.SerializeHashContent(acc.Address)
	return fmt.Sprintf("Hash: %x, Address: %x, TxCnt: %v, Balance: %v, isCreated: %v, isRoot: %v", addressHash[:8], acc.Address[:8], acc.TxCnt, acc.Balance, acc.IsCreated, acc.IsRoot)
}
