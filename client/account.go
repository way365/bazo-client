package client

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/julwil/bazo-client/args"
	"github.com/julwil/bazo-client/cstorage"
	"github.com/julwil/bazo-client/network"
	"github.com/julwil/bazo-client/util"
	"github.com/julwil/bazo-miner/crypto"
	"github.com/julwil/bazo-miner/miner"
	"github.com/julwil/bazo-miner/p2p"
	"github.com/julwil/bazo-miner/protocol"
	"log"
	"math/big"
	"os"
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

func PrepareSignSubmitCreateAccTx(args *args.CreateAccountArgs, logger *log.Logger) error {
	err := args.ValidateInput()
	if err != nil {
		return err
	}

	rootPrivKey, err := crypto.ExtractECDSAKeyFromFile(args.RootWallet)
	if err != nil {
		return err
	}

	// Private Key
	var newPrivKey *ecdsa.PrivateKey

	chParams, err := crypto.GetOrCreateChParamsFromFile(args.ChParams)
	if err != nil {
		return err
	}

	// IMPORTANT: We need to sanitize the secret trapdoor key before we send it to the network.
	chParams.TK = []byte{}

	chCheckString := crypto.NewChCheckString(chParams)

	tx, newPrivKey, err := protocol.ConstrAccTx(
		byte(args.Header),
		uint64(args.Fee),
		[64]byte{},
		rootPrivKey,
		nil,
		nil,
		chParams,
		chCheckString,
		[]byte(args.Data),
	)
	if err != nil {
		return err
	}

	//accAddress := protocol.SerializeHashContent(tx.PubKey)
	//crypto.ChamHashParamsMap[accAddress] = chParams

	//Write the private key to the given textfile
	file, err := os.Create(args.Wallet)
	if err != nil {
		return err
	}

	_, err = file.WriteString(string(newPrivKey.X.Text(16)) + "\n")
	_, err = file.WriteString(string(newPrivKey.Y.Text(16)) + "\n")
	_, err = file.WriteString(string(newPrivKey.D.Text(16)) + "\n")

	if err != nil {
		return errors.New(fmt.Sprintf("failed to write key to file %v", args.Wallet))
	}

	return SendAccountTx(tx, chParams, logger)
}

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

func SendAccountTx(tx protocol.Transaction, chParams *crypto.ChameleonHashParameters, logger *log.Logger) error {

	if err := network.SendTx(util.Config.BootstrapIpport, tx, p2p.ACCTX_BRDCST); err != nil {
		logger.Printf("%v\n", err)
		return err
	} else {
		logger.Printf("Transaction successfully sent to network:\nTxHash: %x%v", tx.ChameleonHash(chParams), tx)
		cstorage.WriteTransaction(tx.ChameleonHash(chParams), tx)
	}

	return nil
}

func AddAccount(args *args.AddAccountArgs, logger *log.Logger) error {
	err := args.ValidateInput()
	if err != nil {
		return err
	}

	rootPrivKey, err := crypto.ExtractECDSAKeyFromFile(args.RootWallet)
	if err != nil {
		return err
	}

	chParams, err := crypto.GetOrCreateChParamsFromFile(args.ChParams)
	if err != nil {
		return err
	}

	chCheckString := crypto.NewChCheckString(chParams)

	var newAddress [64]byte
	newPubInt, _ := new(big.Int).SetString(args.Address, 16)
	copy(newAddress[:], newPubInt.Bytes())

	tx, _, err := protocol.ConstrAccTx(
		byte(args.Header),
		uint64(args.Fee),
		newAddress,
		rootPrivKey,
		nil,
		nil,
		chParams,
		chCheckString,
		[]byte{},
	)
	if err != nil {
		return err
	}

	return SendAccountTx(tx, chParams, logger)
}

func (acc Account) String() string {
	addressHash := protocol.SerializeHashContent(acc.Address)
	return fmt.Sprintf("Hash: %x, Address: %x, TxCnt: %v, Balance: %v, isCreated: %v, isRoot: %v", addressHash[:8], acc.Address[:8], acc.TxCnt, acc.Balance, acc.IsCreated, acc.IsRoot)
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
