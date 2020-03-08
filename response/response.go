package response

import (
	"encoding/json"
	errs "errors"

	"github.com/bob620/baka-rpc-go/errors"
)

type Types string

const (
	SuccessType Types = "success"
	ErrorType   Types = "error"
)

type Response struct {
	responseType Types
	id           string
	jsonRpc      string
	result       *json.RawMessage
	error        *errors.RPCError
}

func NewSuccessResponse(id string, result json.RawMessage) *Response {
	return &Response{
		responseType: SuccessType,
		id:           id,
		jsonRpc:      "2.0",
		result:       &result,
	}
}

func NewErrorResponse(id string, error *errors.RPCError) *Response {
	return &Response{
		responseType: ErrorType,
		id:           id,
		jsonRpc:      "2.0",
		error:        error,
	}
}

func (res *Response) GetRpcVersion() string {
	return res.jsonRpc
}

func (res *Response) GetId() string {
	return res.id
}

func (res *Response) GetType() Types {
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

	// Requires error or result but not both
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

	// Required
	data["jsonrpc"] = []byte(`"` + res.jsonRpc + `"`)

	// May be omitted if error
	if res.id == "" {
		data["id"] = []byte(`null`)
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

	res.jsonRpc = ""
	res.responseType = "error"
	err = json.Unmarshal(jsonData, &jsonReq)
	if err != nil {
		return
	}

	// May be omitted for broken requests
	if jsonReq["id"] != nil {
		var item string
		err = json.Unmarshal(jsonReq["id"], &item)
		if err != nil {
			return
		}

		if item != "" {
			res.id = item
		}
	}

	// Requires error or result but not both
	if jsonReq["error"] != nil {
		err = json.Unmarshal(jsonReq["error"], &res.error)
	} else if jsonReq["result"] != nil {
		res.responseType = SuccessType
		err = json.Unmarshal(jsonReq["result"], &res.result)
	} else {
		return errs.New("no error or result")
	}

	// Don't set the rpc version if an error occurred
	if err != nil {
		return
	}
	err = json.Unmarshal(jsonReq["jsonrpc"], &res.jsonRpc)

	return
}
