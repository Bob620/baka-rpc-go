package rpc

import "encoding/json"

type StringParam struct {
	Name string
	data json.RawMessage
}

func (param *StringParam) getName() string {
	return param.Name
}

func (param *StringParam) getData() json.RawMessage {
	return param.data
}

func (param *StringParam) GetString() (err error, value string) {
	err = json.Unmarshal(param.getData(), &value)
	return
}
