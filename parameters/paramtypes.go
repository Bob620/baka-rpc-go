package parameters

import (
	"encoding/json"
)

type Param interface {
	Clone(json.RawMessage) (Param, error)
	IsRequired() bool
	SetName(string)
	GetName() string
	SetData(json.RawMessage) error
	GetData() json.RawMessage
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
}

type GenericParam struct {
	Name     string
	Default  json.RawMessage
	Required bool
	data     json.RawMessage
}

func (param *GenericParam) Clone(data json.RawMessage) (Param, error) {
	clone := GenericParam{param.Name, param.Default, param.Required, param.data}
	if data != nil {
		err := clone.SetData(data)
		if err != nil {
			return nil, err
		}
	}

	return &clone, nil
}

func (param *GenericParam) IsRequired() bool {
	return param.Required
}

func (param *GenericParam) SetName(newName string) {
	param.Name = newName
}

func (param *GenericParam) GetName() string {
	return param.Name
}

func (param *GenericParam) SetData(message json.RawMessage) (err error) {
	param.data = message
	return
}

func (param *GenericParam) GetData() json.RawMessage {
	return param.data
}

func (param *GenericParam) MarshalJSON() ([]byte, error) {
	if param.data == nil {
		return param.Default, nil
	}
	return param.data, nil
}

func (param *GenericParam) UnmarshalJSON(jsonData []byte) (err error) {
	data := param.Default
	err = json.Unmarshal(jsonData, data)
	return
}

type StringParam struct {
	Name     string
	Default  string
	Required bool
	data     json.RawMessage
}

func (param *StringParam) Clone(data json.RawMessage) (Param, error) {
	clone := StringParam{param.Name, param.Default, param.Required, param.data}
	if data != nil {
		err := clone.SetData(data)
		if err != nil {
			return nil, err
		}
	}

	return &clone, nil
}

func (param *StringParam) IsRequired() bool {
	return param.Required
}

func (param *StringParam) SetName(newName string) {
	param.Name = newName
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

func (param *StringParam) MarshalJSON() ([]byte, error) {
	if param.data == nil {
		data, err := json.Marshal(param.Default)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	return param.data, nil
}

func (param *StringParam) UnmarshalJSON(jsonData []byte) (err error) {
	data := param.Default
	err = json.Unmarshal(jsonData, data)
	return
}

type IntParam struct {
	Name     string
	Default  int
	Required bool
	data     json.RawMessage
}

func (param *IntParam) Clone(data json.RawMessage) (Param, error) {
	clone := IntParam{param.Name, param.Default, param.Required, param.data}
	if data != nil {
		err := clone.SetData(data)
		if err != nil {
			return nil, err
		}
	}

	return &clone, nil
}

func (param *IntParam) IsRequired() bool {
	return param.Required
}

func (param *IntParam) SetName(newName string) {
	param.Name = newName
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

func (param *IntParam) MarshalJSON() ([]byte, error) {
	if param.data == nil {
		data, err := json.Marshal(param.Default)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	return param.data, nil
}

func (param *IntParam) UnmarshalJSON(jsonData []byte) (err error) {
	data := param.Default
	err = json.Unmarshal(jsonData, data)
	return
}

type BoolParam struct {
	Name     string
	Default  bool
	Required bool
	data     json.RawMessage
}

func (param *BoolParam) Clone(data json.RawMessage) (Param, error) {
	clone := BoolParam{param.Name, param.Default, param.Required, param.data}
	if data != nil {
		err := clone.SetData(data)
		if err != nil {
			return nil, err
		}
	}

	return &clone, nil
}

func (param *BoolParam) IsRequired() bool {
	return param.Required
}

func (param *BoolParam) SetName(newName string) {
	param.Name = newName
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

func (param *BoolParam) MarshalJSON() ([]byte, error) {
	if param.data == nil {
		data, err := json.Marshal(param.Default)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	return param.data, nil
}

func (param *BoolParam) UnmarshalJSON(jsonData []byte) (err error) {
	data := param.Default
	err = json.Unmarshal(jsonData, data)
	return
}

type float64Param struct {
	Name     string
	Default  float64
	Required bool
	data     json.RawMessage
}

func (param *float64Param) Clone(data json.RawMessage) (Param, error) {
	clone := float64Param{param.Name, param.Default, param.Required, param.data}
	if data != nil {
		err := clone.SetData(data)
		if err != nil {
			return nil, err
		}
	}

	return &clone, nil
}

func (param *float64Param) IsRequired() bool {
	return param.Required
}

func (param *float64Param) SetName(newName string) {
	param.Name = newName
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

func (param *float64Param) MarshalJSON() ([]byte, error) {
	if param.data == nil {
		data, err := json.Marshal(param.Default)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	return param.data, nil
}

func (param *float64Param) UnmarshalJSON(jsonData []byte) (err error) {
	data := param.Default
	err = json.Unmarshal(jsonData, data)
	return
}
