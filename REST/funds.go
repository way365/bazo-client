package REST

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julwil/bazo-client/args"
	"github.com/julwil/bazo-client/client"
	"github.com/julwil/bazo-client/cstorage"
	"github.com/julwil/bazo-miner/protocol"
	"net/http"
)

type SignatureResponseBody struct {
	TxHash    string `json:"hash"`
	Signature string `json:"signature"`
}

func CreateFundsTx(w http.ResponseWriter, req *http.Request) {
	logger.Println("Incoming createFunds request")
	decoder := json.NewDecoder(req.Body)
	var fundsArgs args.FundsArgs

	err := decoder.Decode(&fundsArgs)
	if err != nil {
		panic(err)
	}

	if fundsArgs.Fee == 0 {
		fundsArgs.Fee = 1
	}

	txHash, err := client.PrepareFundsTx(&fundsArgs, logger)
	if err != nil {
		panic(err)
	}

	var responseBody []Content
	var txResponse Content
	txResponse.Name = "FundsTx"
	txResponse.Detail = fmt.Sprintf("%x", txHash)
	responseBody = append(responseBody, txResponse)

	SendJsonResponse(w, JsonResponse{http.StatusOK, "FundsTx successfully created. Sign the provided hash.", responseBody})
}

func SignFundsTx(w http.ResponseWriter, req *http.Request) {
	logger.Println("Incoming signFunds request")
	decoder := json.NewDecoder(req.Body)
	var requestBody SignatureResponseBody

	err := decoder.Decode(&requestBody)
	if err != nil {
		panic(err)
	}

	var txHash [32]byte
	var Sig [64]byte
	txHashBytes, err := hex.DecodeString(requestBody.TxHash)
	if err != nil {
		panic(err)
	}

	signatureBytes, err := hex.DecodeString(requestBody.Signature)
	if err != nil {
		panic(err)
	}

	copy(txHash[:], txHashBytes)
	copy(Sig[:], signatureBytes)

	tx := cstorage.ReadTransaction(txHash)
	if tx == nil {
		panic(errors.New("transaction not found"))
	}

	tx.SetSignature(Sig)
	cstorage.WriteTransaction(txHash, tx)

	switch tx.(type) {
	case *protocol.FundsTx:
		client.SubmitFundsTx(tx.(*protocol.FundsTx))
	default:
		panic(errors.New("can't cast transaction to funds transaction"))
	}

	var responseBody []Content
	var txResponse Content
	txResponse.Name = "FundsTx"
	txResponse.Detail = fmt.Sprintf("%s", tx.String())
	responseBody = append(responseBody, txResponse)

	SendJsonResponse(w, JsonResponse{http.StatusOK, "FundsTx successfully sent to network.", responseBody})
}
