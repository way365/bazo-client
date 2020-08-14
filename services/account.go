package services

import (
	"github.com/julwil/bazo-client/args"
	"github.com/julwil/bazo-client/cstorage"
	"github.com/julwil/bazo-client/network"
	"github.com/julwil/bazo-miner/crypto"
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

	if err := SignTx(txHash, tx, issuerPrivKey); err != nil {
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

func GetAccount(address [64]byte) (account *protocol.Account, err error) {

	err = network.AccReq(false, protocol.SerializeHashContent(address))
	if err != nil {
		return account, err
	}

	payload, err := network.Fetch(network.AccChan)
	if err != nil {
		return account, err
	}

	account = payload.(*protocol.Account)
	if account.Address == [64]byte{} {
		return account, err
	}

	return account, nil
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
	acc, err := GetAccount(address)
	if err != nil {
		logger.Println(err)
		return err
	} else {
		logger.Printf(acc.String())
	}

	return nil
}
