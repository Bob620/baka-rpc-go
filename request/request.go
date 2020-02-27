package request

import (
	"baka-rpc-go/parameters"
	"encoding/json"
	"github.com/nu7hatch/gouuid"
)

type Request struct {
	requestType string
	id          string
	jsonRpc     string
	method      string
	params      *parameters.Parameters
}

func NewNotification(method string, params *parameters.Parameters) *Request {
	return &Request{
		requestType: "notification",
		jsonRpc:     "2.0",
		method:      method,
		params:      params,
	}
}

func NewRequest(method, id string, params *parameters.Parameters) *Request {
	if len(id) == 0 {
		uid, _ := uuid.NewV4()
		id = uid.String()
	}

	return &Request{
		requestType: "request",
		id:          id,
		jsonRpc:     "2.0",
		method:      method,
		params:      params,
	}
}

func (req *Request) Serialize() (message json.RawMessage, err error) {
	data := map[string]json.RawMessage{}

	params, err := json.Marshal(req.params)

	if err != nil {
		return nil, err
	}

	data["jsonrpc"] = []byte(`"` + req.jsonRpc + `"`)
	data["method"] = []byte(`"` + req.method + `"`)
	data["params"] = params

	if req.requestType != "notification" {
		data["id"] = []byte(`"` + req.id + `"`)
	}

	return json.Marshal(data)
}

func (req *Request) MarshalJSON() ([]byte, error) {
	return req.Serialize()
}

func (req *Request) UnmarshalJSON(jsonData []byte) (err error) {
	var jsonReq map[string]json.RawMessage
	var item string

	req.requestType = "notification"
	err = json.Unmarshal(jsonData, &jsonReq)
	if err != nil {
		return err
	}

	if jsonReq["id"] != nil {
		err = json.Unmarshal(jsonReq["id"], &item)
		if err != nil {
			return err
		}
	}

	if item != "" {
		req.id = item
		req.requestType = "request"
	}

	err = json.Unmarshal(jsonReq["params"], req.params)
	err = json.Unmarshal(jsonReq["method"], &req.method)
	err = json.Unmarshal(jsonReq["jsonrpc"], &req.jsonRpc)

	return
}
