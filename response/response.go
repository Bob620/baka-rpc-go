package response

import (
	"baka-rpc-go/errors"
	"encoding/json"
)

type Response struct {
	responseType string
	id           string
	jsonRpc      string
	result       *json.RawMessage
	error        *errors.RPCError
}

func NewSuccessResponse(id string, result json.RawMessage) *Response {
	return &Response{
		responseType: "success",
		id:           id,
		jsonRpc:      "2.0",
		result:       &result,
	}
}

func NewErrorResponse(id string, error *errors.RPCError) *Response {
	return &Response{
		responseType: "error",
		id:           id,
		jsonRpc:      "2.0",
		error:        error,
	}
}

func (res *Response) GetType() string {
	return res.responseType
}

func (res *Response) GetResult() *json.RawMessage {
	return res.result
}

func (res *Response) GetError() *errors.RPCError {
	return res.error
}

func (res *Response) Serialize() (json.RawMessage, error) {
	data := map[string]json.RawMessage{}

	if res.responseType == "error" {
		resErr, err := json.Marshal(res.error)
		if err != nil {
			return nil, err
		}
		data["error"] = resErr
	} else {
		result, err := json.Marshal(res.result)
		if err != nil {
			return nil, err
		}
		data["result"] = result
	}

	data["jsonrpc"] = []byte(`"` + res.jsonRpc + `"`)

	if res.id == "" {
		data["id"] = []byte(`nil`)
	} else {
		data["id"] = []byte(`"` + res.id + `"`)
	}

	return json.Marshal(data)
}

func (res *Response) MarshalJSON() ([]byte, error) {
	return res.Serialize()
}

func (res *Response) UnmarshalJSON(jsonData []byte) (err error) {
	var jsonReq map[string]json.RawMessage
	var item string

	res.responseType = "error"
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
		res.id = item
	}

	err = json.Unmarshal(jsonReq["jsonrpc"], &res.jsonRpc)
	if jsonReq["error"] != nil {
		err = json.Unmarshal(jsonReq["error"], res.error)
	} else if jsonReq["result"] != nil {
		err = json.Unmarshal(jsonReq["result"], res.result)
	}

	return err
}
