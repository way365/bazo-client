package http

import (
	"encoding/json"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/julwil/bazo-client/util"
	"log"
	"net/http"
)

var (
	logger *log.Logger
)

type SignatureResponseBody struct {
	TxHash    string `json:"hash"`
	Signature string `json:"signature"`
}

func Init() {
	logger = util.InitLogger()

	logger.Printf("%v\n\n", "Starting rest...")

	router := mux.NewRouter()
	getEndpoints(router)
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	ignoreOptions := handlers.IgnoreOptions()

	log.Fatal(http.ListenAndServe(":"+util.Config.Thisclient.Port, handlers.CORS(methodsOk, ignoreOptions)(router)))
}

func getEndpoints(router *mux.Router) {
	// Old code. Needs to be updated.
	//router.HandleFunc("/account/{id}", GetAccountEndpoint).Methods("GET")
	//
	//router.HandleFunc("/createAccTx/{header}/{fee}/{issuer}", CreateAccTxEndpoint).Methods("POST")
	//router.HandleFunc("/createAccTx/{pubKey}/{header}/{fee}/{issuer}", CreateAccTxEndpointWithPubKey).Methods("POST")
	//router.HandleFunc("/sendAccTx/{txHash}/{txSign}", SendAccTxEndpoint).Methods("POST")
	//
	//router.HandleFunc("/createConfigTx/{header}/{id}/{payload}/{fee}/{txCnt}", CreateConfigTxEndpoint).Methods("POST")
	//router.HandleFunc("/sendConfigTx/{txHash}/{txSign}", SendConfigTxEndpoint).Methods("POST")

	router.HandleFunc("/tx/acc", PostAccountTx).Methods("POST")
	router.HandleFunc("/tx/funds", PostFundsTx).Methods("POST")
	router.HandleFunc("/tx/update", PostUpdateTx).Methods("POST")
	router.HandleFunc("/tx/signature", PostSignTx).Methods("POST")

}

func SendJsonResponse(w http.ResponseWriter, resp interface{}) {
	js, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(js)
}
