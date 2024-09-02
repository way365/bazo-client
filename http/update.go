package http

import (
	"encoding/json"
	"fmt"
	"github.com/way365/bazo-client/args"
	"github.com/way365/bazo-client/services"
	"net/http"
)

func PostUpdateTx(w http.ResponseWriter, req *http.Request) {
	logger.Println("Incoming create update request")
	decoder := json.NewDecoder(req.Body)
	var updateTxArgs args.UpdateTxArgs

	err := decoder.Decode(&updateTxArgs)
	if err != nil {
		panic(err)
	}

	err = updateTxArgs.ValidateInput()
	if err != nil {
		fmt.Printf("%v", err)
		SendJsonResponse(w, JsonResponse{http.StatusInternalServerError, "Invalid arguments", []Content{}})
		return
	}

	txHash, _, err := services.PrepareUpdateTx(&updateTxArgs, logger)
	if err != nil {
		panic(err)
	}

	var responseBody []Content
	var txResponse Content
	txResponse.Name = "UpdateTx"
	txResponse.Detail = fmt.Sprintf("%x", txHash)
	responseBody = append(responseBody, txResponse)

	SendJsonResponse(w, JsonResponse{http.StatusOK, "UpdateTx successfully created. Sign the provided hash.", responseBody})
}
