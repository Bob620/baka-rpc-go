package rpc

import (
	"bufio"
	"encoding/json"
	"io"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	UUID "github.com/nu7hatch/gouuid"

	"github.com/bob620/baka-rpc-go/errors"
	"github.com/bob620/baka-rpc-go/parameters"
	"github.com/bob620/baka-rpc-go/request"
	"github.com/bob620/baka-rpc-go/response"
)

type MethodFunc func(params map[string]parameters.Param) (returnMessage json.RawMessage, err error)

type BakaRpc struct {
	chansIn          map[*UUID.UUID]<-chan []byte
	chansOut         map[*UUID.UUID]chan<- []byte
	methods          map[string]*method
	callbackChans    map[string]*chan response.Response
	callbackMutex    sync.RWMutex
	disconnectHandle func(uuid *UUID.UUID)
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
		readerChan <- nil
		_ = conn.Close()
	}()

	return
}

func MakeSocketWriterChan(conn *websocket.Conn) (writerChan chan []byte) {
	writerChan = make(chan []byte)
	go func() {
		evacuate := false
		for !evacuate {
			data := <-writerChan
			if data == nil {
				evacuate = true
			}

			err := conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				evacuate = true
			}
		}
		_ = conn.Close()
	}()

	return
}

func CreateBakaRpc(chanIn <-chan []byte, chanOut chan<- []byte) *BakaRpc {
	rpc := &BakaRpc{
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

func (rpc *BakaRpc) HandleDisconnect(handle func(uuid *UUID.UUID)) {
	rpc.disconnectHandle = handle
}

func (rpc *BakaRpc) AddChannels(chanIn <-chan []byte, chanOut chan<- []byte) (uuid *UUID.UUID) {
	uuid, _ = UUID.NewV4()

	rpc.chansIn[uuid] = chanIn
	rpc.chansOut[uuid] = chanOut

	go rpc.start(uuid)

	return
}

func (rpc *BakaRpc) UseChannels(chanIn <-chan []byte, chanOut chan<- []byte) {
	uuid, _ := UUID.NewV4()

	rpc.chansIn[uuid] = chanIn
	rpc.chansOut[uuid] = chanOut

	rpc.start(uuid)
	rpc.RemoveChannels(uuid)

	return
}

func (rpc *BakaRpc) RemoveChannels(uuid *UUID.UUID) {
	if uuid != nil {
		delete(rpc.chansIn, uuid)
		delete(rpc.chansOut, uuid)
	} else {
		rpc.chansIn = map[*UUID.UUID]<-chan []byte{}
		rpc.chansOut = map[*UUID.UUID]chan<- []byte{}
	}
}

func (rpc *BakaRpc) handleRequest(req request.Request) (message json.RawMessage, errRpc *errors.RPCError) {
	method := rpc.methods[req.GetMethod()]
	if method == nil {
		return nil, errors.NewMethodNotFound()
	}

	numRequiredParams := 0
	reqParams := req.GetParams()
	sanitizedParams := map[string]parameters.Param{}
	// Set default values for everything, assuming we will get all required params
	for _, param := range method.params {
		name := param.GetName()

		// Clone and add if required
		sanitizedParams[name], _ = param.Clone(nil)
		if param.IsRequired() {
			numRequiredParams++
		}

		// Make sure we can assume we get all required params
		if reqParams.Length() < numRequiredParams {
			return nil, errors.NewInvalidParams()
		}
	}

	// Iterate through the request params to change from default values
	switch reqParams.GetType() {
	case parameters.ByName:
		for _, param := range sanitizedParams {
			name := param.GetName()
			reqParam := reqParams.Get(name)

			if reqParam != nil {
				_ = sanitizedParams[name].SetData(reqParam.GetData())
			} else {
				if param.IsRequired() {
					return nil, errors.NewInvalidParams()
				}
			}
		}
		break
	case parameters.ByPosition:
		for index, param := range method.params {
			reqParam := reqParams.Get(strconv.Itoa(index))

			if reqParam != nil {
				_ = sanitizedParams[param.GetName()].SetData(reqParam.GetData())
			} else {
				if param.IsRequired() {
					return nil, errors.NewInvalidParams()
				}
			}
		}
		break
	}

	data, err := (*method.methodFunc)(sanitizedParams)
	if err != nil {
		return nil, errors.NewGenericError(err.Error())
	}

	return data, nil
}

func (rpc *BakaRpc) handleResponse(res response.Response) {
	rpc.callbackMutex.RLock()
	callback := rpc.callbackChans[res.GetId()]
	rpc.callbackMutex.RUnlock()

	if callback != nil {
		*callback <- res
	}

	return
}

func (rpc *BakaRpc) CallMethodByName(channelUuid *UUID.UUID, methodName string, params ...parameters.Param) (res *json.RawMessage, resErr *errors.RPCError) {
	return rpc.CallMethod(channelUuid, methodName, parameters.NewParametersByName(params))
}

func (rpc *BakaRpc) CallMethodByPosition(channelUuid *UUID.UUID, methodName string, params ...parameters.Param) (res *json.RawMessage, resErr *errors.RPCError) {
	return rpc.CallMethod(channelUuid, methodName, parameters.NewParametersByPosition(params))
}

func (rpc *BakaRpc) CallMethodWithNone(channelUuid *UUID.UUID, methodName string) (res *json.RawMessage, resErr *errors.RPCError) {
	return rpc.CallMethod(channelUuid, methodName, &parameters.Parameters{})
}

func (rpc *BakaRpc) CallMethod(channelUuid *UUID.UUID, methodName string, params *parameters.Parameters) (res *json.RawMessage, resErr *errors.RPCError) {
	method := request.NewRequest(methodName, "", params)

	data, err := json.Marshal(method)
	if err != nil {
		return nil, errors.NewParseError()
	}

	callback := make(chan response.Response)
	rpc.callbackMutex.Lock()
	rpc.callbackChans[method.GetId()] = &callback
	rpc.callbackMutex.Unlock()

	if channelUuid == nil {
		for uuid, _ := range rpc.chansOut {
			channelUuid = uuid
			break
		}
	}

	if channelUuid != nil {
		go rpc.sendMessage(data, channelUuid)
		remoteRes := <-callback
		rpc.callbackMutex.Lock()
		delete(rpc.callbackChans, method.GetId())
		rpc.callbackMutex.Unlock()

		if remoteRes.GetType() == response.ErrorType {
			resErr = remoteRes.GetError()
			return nil, resErr
		} else {
			res = remoteRes.GetResult()
			return res, nil
		}
	}
	return nil, errors.NewGenericError("Channel Closed")
}

func (rpc *BakaRpc) NotifyMethodByName(channelUuid *UUID.UUID, methodName string, params ...parameters.Param) {
	rpc.NotifyMethod(channelUuid, methodName, parameters.NewParametersByName(params))
}

func (rpc *BakaRpc) NotifyMethodByPosition(channelUuid *UUID.UUID, methodName string, params ...parameters.Param) {
	rpc.NotifyMethod(channelUuid, methodName, parameters.NewParametersByPosition(params))
}

func (rpc *BakaRpc) NotifyMethodWithNone(channelUuid *UUID.UUID, methodName string) {
	rpc.NotifyMethod(channelUuid, methodName, &parameters.Parameters{})
}

func (rpc *BakaRpc) NotifyMethod(channelUuid *UUID.UUID, methodName string, params *parameters.Parameters) {
	if channelUuid == nil {
		for uuid, _ := range rpc.chansOut {
			channelUuid = uuid
			break
		}
	}

	data, err := json.Marshal(request.NewNotification(methodName, params))
	if err == nil && channelUuid != nil {
		go rpc.sendMessage(data, channelUuid)
	}
}

func (rpc *BakaRpc) start(uuid *UUID.UUID) {
	res := response.Response{}
	req := request.Request{}

	for rpc.chansIn[uuid] != nil {
		message := <-rpc.chansIn[uuid]

		if message == nil {
			rpc.sendMessage(nil, uuid)
			rpc.RemoveChannels(uuid)
			if rpc.disconnectHandle != nil {
				rpc.disconnectHandle(uuid)
			}
			break
		}

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

func (rpc *BakaRpc) sendMessage(message json.RawMessage, uuid *UUID.UUID) {
	rpc.chansOut[uuid] <- message
}

func (rpc *BakaRpc) RegisterMethod(methodName string, methodParams []parameters.Param, methodFunc MethodFunc) {
	rpc.methods[methodName] = &method{
		name:       methodName,
		params:     methodParams,
		methodFunc: &methodFunc,
	}
}

func (rpc *BakaRpc) DeregisterMethod(methodName string) {
	delete(rpc.methods, methodName)
}
