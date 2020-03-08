package rpc

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"strconv"

	"github.com/bob620/baka-rpc-go/errors"
	"github.com/bob620/baka-rpc-go/parameters"
	"github.com/bob620/baka-rpc-go/request"
	"github.com/bob620/baka-rpc-go/response"
)

type MethodFunc func(params map[string]MethodParam) (returnMessage json.RawMessage, err error)

type bakaRpc struct {
	chanIn        <-chan []byte
	chanOut       chan<- []byte
	methods       map[string]*method
	callbackChans map[string]*chan response.Response
}

type method struct {
	name       string
	params     []MethodParam
	methodFunc *MethodFunc
}

type MethodParam interface {
	Clone(json.RawMessage) (MethodParam, error)
	GetName() string
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
		chanIn:        chanIn,
		chanOut:       chanOut,
		methods:       map[string]*method{},
		callbackChans: map[string]*chan response.Response{},
	}
	rpc.start()

	return rpc
}

func (rpc *bakaRpc) handleRequest(req request.Request) (message json.RawMessage, errRpc *errors.RPCError) {
	method := rpc.methods[req.GetMethod()]
	if method == nil {
		return nil, errors.NewMethodNotFound()
	}

	sanitizedParams := map[string]MethodParam{}
	params := req.GetParams()

	if params != nil {
		if params.GetType() == parameters.ByName {
			for _, param := range method.params {
				newParam, err := param.Clone(params.Get(param.GetName()))
				if err != nil {
					return nil, errors.NewInvalidParams()
				}
				sanitizedParams[param.GetName()] = newParam
			}
		} else {
			for key, param := range method.params {
				newParam, err := param.Clone(params.Get(strconv.Itoa(key)))
				if err != nil {
					return nil, errors.NewInvalidParams()
				}
				sanitizedParams[param.GetName()] = newParam
			}
		}
	} else {
		for _, param := range method.params {
			newParam, err := param.Clone(nil)
			if err != nil {
				return nil, errors.NewInvalidParams()
			}
			sanitizedParams[param.GetName()] = newParam
		}
	}

	data, err := (*method.methodFunc)(sanitizedParams)
	if err != nil {
		return nil, errors.NewGenericError("Method failed")
	}

	return data, nil
}

func (rpc *bakaRpc) handleResponse(res response.Response) {
	callback := rpc.callbackChans[res.GetId()]

	if callback != nil {
		*callback <- res
	}

	return
}

func (rpc *bakaRpc) CallMethod(methodName string, params parameters.Parameters) (res *json.RawMessage, resErr *errors.RPCError) {
	method := request.NewRequest(methodName, "", &params)

	data, err := json.Marshal(method)
	if err != nil {
		return nil, errors.NewParseError()
	}

	callback := make(chan response.Response)
	rpc.callbackChans[method.GetId()] = &callback

	go rpc.sendMessage(data)
	remoteRes := <-callback
	delete(rpc.callbackChans, method.GetId())

	if remoteRes.GetType() == response.ErrorType {
		resErr = remoteRes.GetError()
	} else {
		res = remoteRes.GetResult()
	}

	return
}

func (rpc *bakaRpc) NotifyMethod(methodName string, params parameters.Parameters) {
	data, err := json.Marshal(request.NewNotification(methodName, &params))
	if err == nil {
		go rpc.sendMessage(data)
	}
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
					data, _ := json.Marshal(response.NewErrorResponse(res.GetId(), errors.NewParseError()))
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
							data, _ := json.Marshal(response.NewErrorResponse(req.GetId(), err))
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
					go rpc.handleResponse(res)
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

func (rpc *bakaRpc) DeregisterMethod(methodName string) {
	delete(rpc.methods, methodName)
}
