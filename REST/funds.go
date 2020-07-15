package REST

import (
	"encoding/json"
	"fmt"
	"github.com/julwil/bazo-client/args"
	"github.com/julwil/bazo-client/client"
	"net/http"
)

func CreateFundsTxEndpoint2(w http.ResponseWriter, req *http.Request) {
	logger.Println("Incoming createFunds request")
	decoder := json.NewDecoder(req.Body)
	var fundsArgs args.FundsArgs

	err := decoder.Decode(&fundsArgs)
	if err != nil {
		panic(err)
	}

	txHash, err := client.SendFunds(&fundsArgs, logger);
	if err != nil {
		panic(err)
	}

	SendJsonResponse(w, JsonResponse{http.StatusOK, "FundsTx successfully sent to network.", []Content{}})
}
