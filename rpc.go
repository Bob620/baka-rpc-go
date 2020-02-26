package main

import (
	"baka-rpc-go/errors"
	"baka-rpc-go/parameters"
	"encoding/json"
)

type BakaRpc struct {
	chanIn  <-chan []byte
	chanOut chan<- []byte
	methods map[string]*Method
}

type Method struct {
	name   string
	params []MethodParam
}

type MethodParam struct {
	name string
}

func (rpc *BakaRpc) RegisterMethod(method Method, methodFunc func(parameters parameters.Parameters) (*json.RawMessage, *errors.RPCError)) {

}

func (rpc *BakaRpc) DeRegisterMethod(method string) {
	delete(rpc.methods, method)
}
