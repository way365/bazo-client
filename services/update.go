package services

import (
	"errors"
	"fmt"
	"github.com/way365/bazo-client/args"
	"github.com/way365/bazo-client/cstorage"
	"github.com/way365/bazo-miner/crypto"
	"github.com/way365/bazo-miner/protocol"
	"log"
	"math/big"
)

func PrepareSignSubmitUpdateTx(arguments *args.UpdateTxArgs, logger *log.Logger) (txHash [32]byte, err error) {
	txHash, tx, err := PrepareUpdateTx(arguments, logger)
	if err != nil {
		logger.Printf("%v\n", err)
		return [32]byte{}, err
	}

	issuerPrivateKey, err := args.ResolvePrivateKey(arguments.TxIssuer)
	if err != nil {
		return [32]byte{}, err
	}

	if err := SignTx(txHash, tx, issuerPrivateKey); err != nil {
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

func PrepareUpdateTx(arguments *args.UpdateTxArgs, logger *log.Logger) (txHash [32]byte, tx *protocol.UpdateTx, err error) {
	err = arguments.ValidateInput()
	if err != nil {
		return [32]byte{}, tx, err
	}

	// First, we read public key from creator from the wallet file
	issuerPublicKey, err := args.ResolvePublicKey(arguments.TxIssuer)
	if err != nil {
		return [32]byte{}, tx, err
	}

	// Then, we retrieve the associated Address from that private key
	issuerAddress := crypto.GetAddressFromPubKey(issuerPublicKey)

	// Then, we parse the hash of the tx that shall be updated.
	var txToUpdateHash [32]byte
	if len(arguments.TxToUpdate) == 64 {
		newPubInt, _ := new(big.Int).SetString(arguments.TxToUpdate, 16)
		copy(txToUpdateHash[:], newPubInt.Bytes())
	}

	parameters, err := args.ResolveParameters(arguments.Parameters)
	checkString := crypto.NewCheckString(parameters)
	if err != nil {
		return [32]byte{}, tx, errors.New("no chameleon hash parameter files found with given parameters")
	}

	newData := []byte(arguments.UpdateData)
	// We create a new check string for TxToDelete to create a hash collision using chameleon hashing.
	newCheckString := generateCollisionCheckString(txToUpdateHash, parameters, newData)

	// Finally, we create the update-tx.
	tx, err = protocol.ConstrUpdateTx(
		uint64(arguments.Fee),
		txToUpdateHash,
		newCheckString,
		newData,
		protocol.SerializeHashContent(issuerAddress),
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

func generateCollisionCheckString(
	txToUpdateHash [32]byte,
	parameters *crypto.ChameleonHashParameters,
	newData []byte,
) (newCheckString *crypto.ChameleonHashCheckString) {
	// First we need to query the Tx to update.
	var txToUpdate protocol.Transaction
	txToUpdate = cstorage.ReadTransaction(txToUpdateHash)
	if txToUpdate == nil {
		fmt.Printf("TX not found: %x", txToUpdateHash)

		return
	}

	fmt.Printf("TX to update %s", txToUpdate.String())

	// Then we have to save the old check string and the SHA3 hash before we mutate the tx.
	oldCheckString := txToUpdate.GetCheckString()
	oldSHA3 := txToUpdate.SHA3()
	oldHashInput := oldSHA3[:]

	// Now it's time to mutate the tx Data.
	txToUpdate.SetData(newData)

	// Then we compute the new SHA3 hash. This hash incorporates the changes in the Data field.
	// With the new hash input we compute a hash collision and get the new check string.
	newSHA3 := txToUpdate.SHA3()
	newHashInput := newSHA3[:]
	newCheckString = crypto.GenerateChCollision(parameters, oldCheckString, &oldHashInput, &newHashInput)

	// We update the TxToUpdate record in our local db.
	txToUpdate.SetCheckString(newCheckString)
	cstorage.WriteTransaction(txToUpdateHash, txToUpdate)

	return newCheckString
}
