package rpc

import "encoding/json"

type StringParam struct {
	Name    string
	Default string
	data    json.RawMessage
}

func (param *StringParam) SetData(message json.RawMessage) (err error) {
	param.data = message
	_, err = param.GetString()
	return
}

func (param *StringParam) GetData() json.RawMessage {
	return param.data
}

func (param *StringParam) GetString() (value string, err error) {
	err = json.Unmarshal(param.GetData(), &value)
	return
}

type IntParam struct {
	Name    string
	Default int
	data    json.RawMessage
}

func (param *IntParam) SetData(message json.RawMessage) (err error) {
	param.data = message
	_, err = param.GetInt()
	return
}

func (param *IntParam) GetData() json.RawMessage {
	return param.data
}

func (param *IntParam) GetInt() (value int, err error) {
	err = json.Unmarshal(param.GetData(), &value)
	return
}

type BoolParam struct {
	Name    string
	Default bool
	data    json.RawMessage
}

func (param *BoolParam) SetData(message json.RawMessage) (err error) {
	param.data = message
	_, err = param.GetBool()
	return
}

func (param *BoolParam) GetData() json.RawMessage {
	return param.data
}

func (param *BoolParam) GetBool() (value bool, err error) {
	err = json.Unmarshal(param.GetData(), &value)
	return
}

type float64Param struct {
	Name    string
	Default float64
	data    json.RawMessage
}

func (param *float64Param) SetData(message json.RawMessage) (err error) {
	param.data = message
	_, err = param.GetFloat64()
	return
}

func (param *float64Param) GetData() json.RawMessage {
	return param.data
}

func (param *float64Param) GetFloat64() (value float64, err error) {
	err = json.Unmarshal(param.GetData(), &value)
	return
}
