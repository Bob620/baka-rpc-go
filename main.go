package main

import (
	"encoding/json"
	"fmt"

	"github.com/bob620/baka-rpc-go/parameters"
	"github.com/bob620/baka-rpc-go/rpc"
)

func main() {
	// Testing Overhead
	chanOne := make(chan []byte)
	chanTwo := make(chan []byte)

	// Client One
	rpcClient := rpc.CreateBakaRpc(chanOne, chanTwo)

	// Register one method
	rpcClient.RegisterMethod(
		"idk",
		[]parameters.Param{
			&parameters.StringParam{Name: "test"},
		}, func(params map[string]parameters.Param) (returnMessage json.RawMessage, err error) {
			test, _ := params["test"].(*parameters.StringParam).GetString()

			return json.Marshal(test)
		})

	// Client Two
	rpcClient2 := rpc.CreateBakaRpc(chanTwo, chanOne)

	// Request one method
	res, resErr := rpcClient2.CallMethod("idk",
		parameters.NewParametersByPosition([]parameters.Param{
			&parameters.StringParam{Name: "test", Default: "ahhh"},
		}))

	// Handle method return
	if resErr != nil {
		fmt.Printf("%s\n", resErr.Message)
		return
	}

	data := ""
	err := json.Unmarshal(*res, &data)
	if err != nil {
		println(err.Error())
		return
	}

	fmt.Printf("Response: `%s`\n", data)
}
