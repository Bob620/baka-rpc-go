package main

import (
	"baka-rpc-go/parameters"
	"encoding/json"
	"fmt"
)

func main() {
	noMethodError := NewMethodNotFound()
	testString := "test"
	output, err := json.Marshal(NewErrorResponse(&testString, noMethodError))
	if err != nil {
		return
	}

	fmt.Printf("%s\n", output)

	result := []byte(`1`)
	output, err = json.Marshal(NewSuccessResponse(&testString, result))
	if err != nil {
		return
	}

	fmt.Printf("%s\n", output)

	params := parameters.NewParametersByPosition()
	params.SetString("3", "\"hi")
	params.SetFloat("2", 0.7)
	output, err = json.Marshal(NewRequest("", "", params))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%s\n", output)

	params = parameters.NewParametersByPosition()
	var test = NewRequest("", "", params)
	err = json.Unmarshal(output, &test)
	if err != nil {
		fmt.Println(err)
		return
	}

	output, err = json.Marshal(test)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%s\n", output)
}
