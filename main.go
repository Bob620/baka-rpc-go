package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/bob620/baka-rpc-go/parameters"
	"github.com/bob620/baka-rpc-go/rpc"
)

// Testing Overhead
var upgrader = websocket.Upgrader{} // use default options

func main() {
	// Client One
	rpcClient := rpc.CreateBakaRpc(nil, nil)

	// Register one method
	rpcClient.RegisterMethod(
		"idk",
		[]parameters.Param{
			&parameters.StringParam{Name: "test"},
		}, func(params map[string]parameters.Param) (returnMessage json.RawMessage, err error) {
			test, _ := params["test"].(*parameters.StringParam).GetString()

			return json.Marshal(test)
		})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		defer c.Close()
		rpcClient.UseChannels(rpc.MakeSocketReaderChan(c), rpc.MakeSocketWriterChan(c))
	})

	http.ListenAndServe("localhost:9889", nil)
}
