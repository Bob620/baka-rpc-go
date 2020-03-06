package rpc

import (
	"encoding/json"
	"errors"
)

type StringParam struct {
	Name         string
	Default      string
	NeedsDefault bool
	data         json.RawMessage
}

func (param *StringParam) Clone(message json.RawMessage) (MethodParam, error) {
	if param.NeedsDefault && message == nil {
		return nil, errors.New("requires non-nil value")
	}
	clone := StringParam{Default: param.Default, Name: param.Name}
	err := clone.SetData(message)
	return &clone, err
}

func (param *StringParam) GetName() string {
	return param.Name
}

func (param *StringParam) SetData(message json.RawMessage) (err error) {
	param.data = message
	_, err = param.GetString()
	return
}

func (param *StringParam) GetData() json.RawMessage {
	if param.data == nil {
		data, _ := json.Marshal(param.Default)
		return data
	}
	return param.data
}

func (param *StringParam) GetString() (value string, err error) {
	err = json.Unmarshal(param.GetData(), &value)
	return
}

type IntParam struct {
	Name         string
	Default      int
	NeedsDefault bool
	data         json.RawMessage
}

func (param *IntParam) Clone(message json.RawMessage) (MethodParam, error) {
	if param.NeedsDefault && message == nil {
		return nil, errors.New("requires non-nil value")
	}
	clone := IntParam{Default: param.Default, Name: param.Name}
	err := clone.SetData(message)
	return &clone, err
}

func (param *IntParam) GetName() string {
	return param.Name
}

func (param *IntParam) SetData(message json.RawMessage) (err error) {
	param.data = message
	_, err = param.GetInt()
	return
}

func (param *IntParam) GetData() json.RawMessage {
	if param.data == nil {
		data, _ := json.Marshal(param.Default)
		return data
	}
	return param.data
}

func (param *IntParam) GetInt() (value int, err error) {
	err = json.Unmarshal(param.GetData(), &value)
	return
}

type BoolParam struct {
	Name         string
	Default      bool
	NeedsDefault bool
	data         json.RawMessage
}

func (param *BoolParam) Clone(message json.RawMessage) (MethodParam, error) {
	if param.NeedsDefault && message == nil {
		return nil, errors.New("requires non-nil value")
	}
	clone := BoolParam{Default: param.Default, Name: param.Name}
	err := clone.SetData(message)
	return &clone, err
}

func (param *BoolParam) GetName() string {
	return param.Name
}

func (param *BoolParam) SetData(message json.RawMessage) (err error) {
	param.data = message
	_, err = param.GetBool()
	return
}

func (param *BoolParam) GetData() json.RawMessage {
	if param.data == nil {
		data, _ := json.Marshal(param.Default)
		return data
	}
	return param.data
}

func (param *BoolParam) GetBool() (value bool, err error) {
	err = json.Unmarshal(param.GetData(), &value)
	return
}

type float64Param struct {
	Name         string
	Default      float64
	NeedsDefault bool
	data         json.RawMessage
}

func (param *float64Param) Clone(message json.RawMessage) (MethodParam, error) {
	if param.NeedsDefault && message == nil {
		return nil, errors.New("requires non-nil value")
	}
	clone := float64Param{Default: param.Default, Name: param.Name}
	err := clone.SetData(message)
	return &clone, err
}

func (param *float64Param) GetName() string {
	return param.Name
}

func (param *float64Param) SetData(message json.RawMessage) (err error) {
	param.data = message
	_, err = param.GetFloat64()
	return
}

func (param *float64Param) GetData() json.RawMessage {
	if param.data == nil {
		data, _ := json.Marshal(param.Default)
		return data
	}
	return param.data
}

func (param *float64Param) GetFloat64() (value float64, err error) {
	err = json.Unmarshal(param.GetData(), &value)
	return
}
