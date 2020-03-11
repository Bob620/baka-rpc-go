package rpc

import (
	"bufio"
	"encoding/json"
	"io"
	"strconv"

	"github.com/gorilla/websocket"
	UUID "github.com/nu7hatch/gouuid"

	"github.com/bob620/baka-rpc-go/errors"
	"github.com/bob620/baka-rpc-go/parameters"
	"github.com/bob620/baka-rpc-go/request"
	"github.com/bob620/baka-rpc-go/response"
)

type MethodFunc func(params map[string]parameters.Param) (returnMessage json.RawMessage, err error)

type bakaRpc struct {
	chansIn       map[*UUID.UUID]<-chan []byte
	chansOut      map[*UUID.UUID]chan<- []byte
	methods       map[string]*method
	callbackChans map[string]*chan response.Response
}

type method struct {
	name       string
	params     []parameters.Param
	methodFunc *MethodFunc
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

func MakeSocketReaderChan(conn *websocket.Conn) (readerChan chan []byte) {
	readerChan = make(chan []byte)
	go func() {
		evacuate := false
		for !evacuate {
			_, message, err := conn.ReadMessage()
			if err != nil {
				evacuate = true
			}
			readerChan <- message
		}
	}()

	return
}

func MakeSocketWriterChan(conn *websocket.Conn) (writerChan chan []byte) {
	writerChan = make(chan []byte)
	go func() {
		evacuate := false
		for !evacuate {
			err := conn.WriteMessage(websocket.TextMessage, <-writerChan)
			if err != nil {
				evacuate = true
			}
		}
	}()

	return
}

func CreateBakaRpc(chanIn <-chan []byte, chanOut chan<- []byte) *bakaRpc {
	rpc := &bakaRpc{
		chansIn:       map[*UUID.UUID]<-chan []byte{},
		chansOut:      map[*UUID.UUID]chan<- []byte{},
		methods:       map[string]*method{},
		callbackChans: map[string]*chan response.Response{},
	}
	if chanIn != nil && chanOut != nil {
		rpc.AddChannels(chanIn, chanOut)
	}

	return rpc
}

func (rpc *bakaRpc) AddChannels(chanIn <-chan []byte, chanOut chan<- []byte) (uuid *UUID.UUID) {
	uuid, _ = UUID.NewV4()

	rpc.chansIn[uuid] = chanIn
	rpc.chansOut[uuid] = chanOut

	go rpc.start(uuid)

	return
}

func (rpc *bakaRpc) UseChannels(chanIn <-chan []byte, chanOut chan<- []byte) {
	uuid, _ := UUID.NewV4()

	rpc.chansIn[uuid] = chanIn
	rpc.chansOut[uuid] = chanOut

	rpc.start(uuid)
	rpc.RemoveChannels(uuid)

	return
}

func (rpc *bakaRpc) RemoveChannels(uuid *UUID.UUID) {
	if uuid != nil {
		delete(rpc.chansIn, uuid)
		delete(rpc.chansOut, uuid)
	}
}

func (rpc *bakaRpc) handleRequest(req request.Request) (message json.RawMessage, errRpc *errors.RPCError) {
	method := rpc.methods[req.GetMethod()]
	if method == nil {
		return nil, errors.NewMethodNotFound()
	}

	sanitizedParams := map[string]parameters.Param{}
	params := req.GetParams()

	if params != nil {
		if params.GetType() == parameters.ByName {
			for _, param := range method.params {
				reqParam := params.Get(param.GetName())
				if reqParam == nil {
					return nil, errors.NewInvalidParams()
				}
				newParam, err := param.Clone(reqParam.GetData())
				if err != nil {
					return nil, errors.NewInvalidParams()
				}
				sanitizedParams[param.GetName()] = newParam
			}
		} else {
			for key, param := range method.params {
				reqParam := params.Get(strconv.Itoa(key))
				if reqParam == nil {
					return nil, errors.NewInvalidParams()
				}
				newParam, err := param.Clone(reqParam.GetData())
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

func (rpc *bakaRpc) CallMethod(channelUuid *UUID.UUID, methodName string, params *parameters.Parameters) (res *json.RawMessage, resErr *errors.RPCError) {
	method := request.NewRequest(methodName, "", params)

	data, err := json.Marshal(method)
	if err != nil {
		return nil, errors.NewParseError()
	}

	callback := make(chan response.Response)
	rpc.callbackChans[method.GetId()] = &callback

	if channelUuid == nil {
		for uuid, _ := range rpc.chansOut {
			channelUuid = uuid
			break
		}
	}

	if channelUuid != nil {
		go rpc.sendMessage(data, channelUuid)
		remoteRes := <-callback
		delete(rpc.callbackChans, method.GetId())

		if remoteRes.GetType() == response.ErrorType {
			resErr = remoteRes.GetError()
		} else {
			res = remoteRes.GetResult()
		}
	}
	return nil, errors.NewGenericError("Channel Closed")
}

func (rpc *bakaRpc) NotifyMethod(channelUuid *UUID.UUID, methodName string, params parameters.Parameters) {
	if channelUuid == nil {
		for uuid, _ := range rpc.chansOut {
			channelUuid = uuid
			break
		}
	}

	data, err := json.Marshal(request.NewNotification(methodName, &params))
	if err == nil && channelUuid != nil {
		go rpc.sendMessage(data, channelUuid)
	}
}

func (rpc *bakaRpc) start(uuid *UUID.UUID) {
	res := response.Response{}
	req := request.Request{}

	for rpc.chansIn[uuid] != nil {
		message := <-rpc.chansIn[uuid]

		err := json.Unmarshal(message, &req)
		if err != nil {
			err = json.Unmarshal(message, &res)
			if err != nil {
				data, _ := json.Marshal(response.NewErrorResponse(res.GetId(), errors.NewParseError()))
				go rpc.sendMessage(data, uuid)
			}
		}

		if req.GetRpcVersion() != "" {
			if req.GetRpcVersion() != "2.0" {
				data, _ := json.Marshal(response.NewErrorResponse(req.GetId(), errors.NewInvalidRequest()))
				go rpc.sendMessage(data, uuid)
			} else {
				go func() {
					message, err := rpc.handleRequest(req)
					if err != nil {
						data, _ := json.Marshal(response.NewErrorResponse(req.GetId(), err))
						rpc.sendMessage(data, uuid)
					} else {
						data, _ := json.Marshal(response.NewSuccessResponse(req.GetId(), message))
						rpc.sendMessage(data, uuid)
					}
				}()
			}
		}

		if res.GetRpcVersion() != "" {
			if res.GetRpcVersion() != "2.0" {
				data, _ := json.Marshal(response.NewErrorResponse(res.GetId(), errors.NewInvalidRequest()))
				go rpc.sendMessage(data, uuid)
			} else {
				go rpc.handleResponse(res)
			}
		}
	}
}

func (rpc *bakaRpc) sendMessage(message json.RawMessage, uuid *UUID.UUID) {
	rpc.chansOut[uuid] <- message
}

func (rpc *bakaRpc) RegisterMethod(methodName string, methodParams []parameters.Param, methodFunc MethodFunc) {
	rpc.methods[methodName] = &method{
		name:       methodName,
		params:     methodParams,
		methodFunc: &methodFunc,
	}
}

func (rpc *bakaRpc) DeregisterMethod(methodName string) {
	delete(rpc.methods, methodName)
}
