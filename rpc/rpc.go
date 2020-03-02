package rpc

import (
	"baka-rpc-go/errors"
	"baka-rpc-go/response"
	"bufio"
	"encoding/json"
	"io"
	"os"
)

type MethodFunc func(params map[string]MethodParam) (returnMessage json.RawMessage, err error)

type bakaRpc struct {
	chanIn  <-chan []byte
	chanOut chan<- []byte
	methods map[string]*method
}

type method struct {
	name       string
	params     []MethodParam
	methodFunc *MethodFunc
}

type MethodParam interface {
	GetData() json.RawMessage
}

func MakeReaderChan(r io.Reader) <-chan []byte {
	data := make(chan []byte)
	go func() {
		defer close(data)
		scan := bufio.NewScanner(r)
		for scan.Scan() {
			data <- scan.Bytes()
		}
	}()
	return data
}

func MakeWriterChan(r io.Writer) chan<- []byte {
	data := make(chan []byte)
	go func() {
		defer close(data)
		write := bufio.NewWriter(r)
		for {
			_, err := write.Write(<-data)
			if err != nil {
				break
			}
		}
	}()
	return data
}

func CreateBakaRpc(chanIn <-chan []byte, chanOut chan<- []byte) *bakaRpc {
	if chanIn == nil {
		chanIn = MakeReaderChan(os.Stdin)
	}

	if chanOut == nil {
		chanOut = MakeWriterChan(os.Stdout)
	}

	return &bakaRpc{
		chanIn:  chanIn,
		chanOut: chanOut,
		methods: map[string]*method{},
	}
}

func (rpc *bakaRpc) RegisterMethod(methodName string, methodParams []MethodParam, methodFunc MethodFunc) {
	rpc.methods[methodName] = &method{
		name:       methodName,
		params:     methodParams,
		methodFunc: &methodFunc,
	}
}

func (rpc *bakaRpc) DeRegisterMethod(methodName string) {
	delete(rpc.methods, methodName)
}

func (rpc *bakaRpc) CallMethod(methodName string, methodParams map[string]MethodParam) *response.Response {
	method := rpc.methods[methodName]
	if method == nil {
		return response.NewErrorResponse("", errors.NewMethodNotFound())
	}

	returnMessage, err := (*method.methodFunc)(methodParams)

	if err != nil {
		return response.NewErrorResponse("", errors.NewGenericError(err.Error()))
	}

	return response.NewSuccessResponse("", returnMessage)
}
