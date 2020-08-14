package http

import (
	"encoding/json"
	"fmt"
	"github.com/julwil/bazo-client/args"
	"github.com/julwil/bazo-client/services"
	"net/http"
)

func PostFundsTx(w http.ResponseWriter, req *http.Request) {
	logger.Println("Incoming createFunds request")
	decoder := json.NewDecoder(req.Body)
	var fundsArgs args.FundsArgs

	err := decoder.Decode(&fundsArgs)
	if err != nil {
		panic(err)
	}

	err = fundsArgs.ValidateInput()
	if err != nil {
		fmt.Printf("%v", err)
		SendJsonResponse(w, JsonResponse{http.StatusInternalServerError, "Invalid arguments", []Content{}})
		return
	}

	txHash, _, err := services.PrepareFundsTx(&fundsArgs, logger)
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
