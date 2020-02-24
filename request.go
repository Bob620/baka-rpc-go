package main

import (
	"baka-rpc-go/parameters"
	"github.com/nu7hatch/gouuid"
)

type Request struct {
	JsonRpc string                 `json:"jsonrpc"`
	Id      string                 `json:"id"`
	Method  string                 `json:"method"`
	Params  *parameters.Parameters `json:"params"`
}

type Notification struct {
	JsonRpc string                 `json:"jsonrpc"`
	Method  string                 `json:"method"`
	Params  *parameters.Parameters `json:"params"`
}

func NewNotification(method string, params *parameters.Parameters) *Notification {
	return &Notification{
		JsonRpc: "2.0",
		Method:  method,
		Params:  params,
	}
}

func NewRequest(method, id string, params *parameters.Parameters) *Request {
	if len(id) == 0 {
		uid, _ := uuid.NewV4()
		id = uid.String()
	}

	return &Request{
		JsonRpc: "2.0",
		Id:      id,
		Method:  method,
		Params:  params,
	}
}
