package main

import (
	"baka-rpc-go/errors"
	"baka-rpc-go/parameters"
	"baka-rpc-go/request"
	"baka-rpc-go/response"
	"encoding/json"
	"fmt"
)

func main() {
	noMethodError := errors.NewMethodNotFound()
	testString := "testReq"
	output, err := json.Marshal(response.NewErrorResponse(testString, noMethodError))
	if err != nil {
		return
	}

	fmt.Printf("%s\n", output)

	result := []byte(`1`)
	output, err = json.Marshal(response.NewSuccessResponse(testString, result))
	if err != nil {
		return
	}

	fmt.Printf("%s\n", output)

	params := parameters.NewParametersByPosition()
	params.SetString("3", "\"hi")
	params.SetFloat("2", 0.7)
	output, err = json.Marshal(request.NewRequest("", "", params))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%s\n", output)

	params = parameters.NewParametersByPosition()
	testReq := request.NewRequest("", "", params)
	err = json.Unmarshal(output, &testReq)
	if err != nil {
		fmt.Println(err)
		return
	}

	output, err = json.Marshal(testReq)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%s\n", output)

	output, err = json.Marshal(request.NewNotification("", params))
	if err != nil {
		fmt.Println(err)
		return
	}

	params = parameters.NewParametersByPosition()
	testNotif := request.NewRequest("", "", params)
	err = json.Unmarshal(output, &testNotif)
	if err != nil {
		fmt.Println(err)
		return
	}

	output, err = json.Marshal(testNotif)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%s\n", output)
}
