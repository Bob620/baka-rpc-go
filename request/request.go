package request

import (
	"baka-rpc-go/parameters"
	"encoding/json"
	"github.com/nu7hatch/gouuid"
)

type Request struct {
	requestType string
	id          string
	JsonRpc     string
	Method      string
	Params      *parameters.Parameters
}

func NewNotification(method string, params *parameters.Parameters) *Request {
	return &Request{
		requestType: "notification",
		JsonRpc:     "2.0",
		Method:      method,
		Params:      params,
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
		JsonRpc:     "2.0",
		Method:      method,
		Params:      params,
	}
}

func (req *Request) Serialize() (json.RawMessage, error) {
	data := map[string]json.RawMessage{}
	var err error

	params, err := json.Marshal(req.Params)

	if err != nil {
		return nil, err
	}

	data["jsonrpc"] = []byte(`"` + req.JsonRpc + `"`)
	data["method"] = []byte(`"` + req.Method + `"`)
	data["params"] = params

	if req.requestType != "notification" {
		data["id"] = []byte(`"` + req.id + `"`)
	}

	return json.Marshal(data)
}

func (req *Request) MarshalJSON() ([]byte, error) {
	return req.Serialize()
}

func (req *Request) UnmarshalJSON(jsonData []byte) error {
	var jsonReq map[string]json.RawMessage
	var item string
	var err error

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

	err = json.Unmarshal(jsonReq["params"], req.Params)
	err = json.Unmarshal(jsonReq["method"], &req.Method)
	err = json.Unmarshal(jsonReq["jsonrpc"], &req.JsonRpc)

	return err
}
