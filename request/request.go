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

func (req *Request) GetRpcVersion() string {
	return req.jsonRpc
}

func (req *Request) GetId() string {
	return req.id
}

func (req *Request) Serialize() (message json.RawMessage, err error) {
	data := map[string]json.RawMessage{}

	params, err := json.Marshal(req.params)
	if err != nil {
		return nil, err
	}

	// Required
	data["jsonrpc"] = []byte(`"` + req.jsonRpc + `"`)
	data["method"] = []byte(`"` + req.method + `"`)

	// May be omitted
	if data["params"] != nil {
		data["params"] = params
	}

	// Omitted if notification
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

	// It's a notification until it gets an id
	req.jsonRpc = ""
	req.requestType = "notification"
	err = json.Unmarshal(jsonData, &jsonReq)
	if err != nil {
		return err
	}

	// May be omitted for Notifications
	if jsonReq["id"] != nil {
		var item string
		err = json.Unmarshal(jsonReq["id"], &item)
		if err != nil {
			return err
		}

		if item != "" {
			req.id = item
			req.requestType = "request"
		}
	}

	// May be omitted
	if jsonReq["params"] != nil {
		err = json.Unmarshal(jsonReq["params"], req.params)
		if err != nil {
			return err
		}
	}

	// Always required
	err = json.Unmarshal(jsonReq["method"], &req.method)

	// Don't set the rpc version if an error occurred
	if err == nil {
		return
	}
	err = json.Unmarshal(jsonReq["jsonrpc"], &req.jsonRpc)

	return
}
