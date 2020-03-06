package main

import (
	"encoding/json"
	"fmt"

	"baka-rpc-go/parameters"
	"baka-rpc-go/rpc"
)

func main() {
	chanOne := make(chan []byte)
	chanTwo := make(chan []byte)

	rpcClient := rpc.CreateBakaRpc(chanOne, chanTwo)
	rpcClient2 := rpc.CreateBakaRpc(chanTwo, chanOne)
	rpcClient.RegisterMethod(
		"idk",
		[]rpc.MethodParam{
			&rpc.StringParam{Name: "test"},
		}, func(params map[string]rpc.MethodParam) (returnMessage json.RawMessage, err error) {
			test, _ := params["test"].(*rpc.StringParam).GetString()

			return json.Marshal(test)
		})

	params := parameters.NewParametersByName()
	params.SetString("test", "ahhhh")
	res, resErr := rpcClient2.CallMethod("idk", *params)
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
