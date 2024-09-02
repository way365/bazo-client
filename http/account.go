package http

import (
	"encoding/json"
	"fmt"
	"github.com/way365/bazo-client/args"
	"github.com/way365/bazo-client/services"
	"net/http"
)

func PostAccountTx(w http.ResponseWriter, req *http.Request) {
	logger.Println("Incoming create account request")
	decoder := json.NewDecoder(req.Body)
	var createAccountArgs args.CreateAccountArgs

	err := decoder.Decode(&createAccountArgs)
	if err != nil {
		panic(err)
	}

	err = createAccountArgs.ValidateInput()
	if err != nil {
		fmt.Printf("%v", err)
		SendJsonResponse(w, JsonResponse{http.StatusInternalServerError, "Invalid arguments", []Content{}})
		return
	}

	txHash, _, err := services.PrepareCreateAccountTx(&createAccountArgs, logger)
	if err != nil {
		panic(err)
	}

	var responseBody []Content
	var txResponse Content
	txResponse.Name = "CreateAccountTx"
	txResponse.Detail = fmt.Sprintf("%x", txHash)
	responseBody = append(responseBody, txResponse)

	SendJsonResponse(w, JsonResponse{http.StatusOK, "AccountTx successfully created. Sign the provided hash.", responseBody})
}
