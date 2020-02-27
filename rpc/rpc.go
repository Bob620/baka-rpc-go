package rpc

import (
	"baka-rpc-go/parameters"
	"encoding/json"
)

type BakaRpc struct {
	chanIn  <-chan []byte
	chanOut chan<- []byte
	methods map[string]*Method
}

type Method struct {
	Name   string
	Params []MethodParam
}

type MethodParam interface {
	getName() string
	getData() json.RawMessage
}

func (rpc *BakaRpc) RegisterMethod(method Method, methodFunc func(params parameters.Parameters) (*json.RawMessage, error)) {

}

func (rpc *BakaRpc) DeRegisterMethod(method string) {
	delete(rpc.methods, method)
}
