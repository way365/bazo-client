package main

import (
	"github.com/mchetelat/bazo_client/REST"
	"github.com/mchetelat/bazo_client/client"
	"os"
)

func main() {
	client.Init()
	if len(os.Args) >= 2 {
		if os.Args[1] == "accTx" || os.Args[1] == "fundsTx" || os.Args[1] == "configTx" || os.Args[1] == "stakeTx" {
			client.Process(os.Args[1:])
		} else {
			client.State(os.Args[1])
		}
	} else {
		REST.Init()
	}
}
