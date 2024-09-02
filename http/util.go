package http

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/way365/bazo-client/cstorage"
	"github.com/way365/bazo-client/services"
	"net/http"
)

func PostSignTx(w http.ResponseWriter, req *http.Request) {
	logger.Println("Incoming sign tx request")
	decoder := json.NewDecoder(req.Body)
	var requestBody SignatureResponseBody

	err := decoder.Decode(&requestBody)
	if err != nil {
		panic(err)
	}

	var txHash [32]byte
	var Signature [64]byte

	signatureBytes, err := hex.DecodeString(requestBody.Signature)
	txHashBytes, err := hex.DecodeString(requestBody.TxHash)
	if err != nil {
		panic(err)
	}

	copy(txHash[:], txHashBytes[:])
	copy(Signature[:], signatureBytes[:])

	tx := cstorage.ReadTransaction(txHash)
	if tx == nil {
		panic(errors.New("transaction not found"))
	}

	tx.SetSignature(Signature)

	services.SubmitTx(txHash, tx)

	cstorage.WriteTransaction(txHash, tx)

	var responseBody []Content
	var txResponse Content
	txResponse.Name = "Transaction"

	txResponse.Detail = fmt.Sprintf("%x", txHash)
	responseBody = append(responseBody, txResponse)

	SendJsonResponse(w, JsonResponse{http.StatusOK, "Tx successfully sent to network.", responseBody})
}
