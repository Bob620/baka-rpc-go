package rpc

import "encoding/json"

type StringParam struct {
	Name    string
	Default string
	data    json.RawMessage
}

func (param StringParam) GetData() json.RawMessage {
	return param.data
}

func (param StringParam) GetString() (value string, err error) {
	err = json.Unmarshal(param.GetData(), &value)
	return
}

type IntParam struct {
	Name    string
	Default int
	data    json.RawMessage
}

func (param IntParam) GetData() json.RawMessage {
	return param.data
}

func (param IntParam) GetInt() (value int, err error) {
	err = json.Unmarshal(param.GetData(), &value)
	return
}

type BoolParam struct {
	Name    string
	Default bool
	data    json.RawMessage
}

func (param BoolParam) GetData() json.RawMessage {
	return param.data
}

func (param BoolParam) GetInt() (value bool, err error) {
	err = json.Unmarshal(param.GetData(), &value)
	return
}

type float64Param struct {
	Name    string
	Default float64
	data    json.RawMessage
}

func (param float64Param) GetData() json.RawMessage {
	return param.data
}

func (param float64Param) GetInt() (value float64, err error) {
	err = json.Unmarshal(param.GetData(), &value)
	return
}
