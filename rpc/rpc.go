package rpc

import (
	"baka-rpc-go/errors"
	"baka-rpc-go/request"
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
	SetData(json.RawMessage) error
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

	rpc := &bakaRpc{
		chanIn:  chanIn,
		chanOut: chanOut,
		methods: map[string]*method{},
	}
	rpc.start()

	return rpc
}

func (rpc *bakaRpc) handleRequest(req request.Request) (message json.RawMessage, err error) {
	return
}

func (rpc *bakaRpc) handleResponse(res response.Response) (message json.RawMessage, err error) {
	return
}

func (rpc *bakaRpc) start() {
	go func() {
		res := response.Response{}
		req := request.Request{}

		for {
			message := <-rpc.chanIn

			err := json.Unmarshal(message, &req)
			if err != nil {
				err = json.Unmarshal(message, &res)
				if err != nil {
					data, _ := json.Marshal(response.NewErrorResponse("", errors.NewParseError()))
					go rpc.sendMessage(data)
				}
			}

			if req.GetRpcVersion() != "" {
				if req.GetRpcVersion() != "2.0" {
					data, _ := json.Marshal(response.NewErrorResponse(req.GetId(), errors.NewInvalidRequest()))
					go rpc.sendMessage(data)
				} else {
					go func() {
						message, err := rpc.handleRequest(req)
						if err != nil {
							data, _ := json.Marshal(response.NewErrorResponse(req.GetId(), errors.NewGenericError(err.Error())))
							rpc.sendMessage(data)
						} else {
							data, _ := json.Marshal(response.NewSuccessResponse(req.GetId(), message))
							rpc.sendMessage(data)
						}
					}()
				}
			}

			if res.GetRpcVersion() != "" {
				if res.GetRpcVersion() != "2.0" {
					data, _ := json.Marshal(response.NewErrorResponse(res.GetId(), errors.NewInvalidRequest()))
					go rpc.sendMessage(data)
				} else {
					go func() {
						message, err := rpc.handleResponse(res)
						if err != nil {
							data, _ := json.Marshal(response.NewErrorResponse(res.GetId(), errors.NewGenericError(err.Error())))
							rpc.sendMessage(data)
						} else {
							data, _ := json.Marshal(response.NewSuccessResponse(res.GetId(), message))
							rpc.sendMessage(data)
						}
					}()
				}
			}
		}
	}()
}

func (rpc *bakaRpc) sendMessage(message json.RawMessage) {
	rpc.chanOut <- message
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

func (rpc *bakaRpc) CallLocalMethod(methodName string, methodParams map[string]MethodParam) (success *response.Response, error *response.Response) {
	method := rpc.methods[methodName]
	if method == nil {
		return nil, response.NewErrorResponse("", errors.NewMethodNotFound())
	}

	returnMessage, err := (*method.methodFunc)(methodParams)

	if err != nil {
		return nil, response.NewErrorResponse("", errors.NewGenericError(err.Error()))
	}

	return response.NewSuccessResponse("", returnMessage), nil
}
