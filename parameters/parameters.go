package parameters

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

type Types string

const (
	ByName     Types = "byName"
	ByPosition Types = "byPosition"
)

type Parameters struct {
	paramType Types
	values    map[string]json.RawMessage
}

func newParameters(paramType Types) *Parameters {
	return &Parameters{paramType: paramType, values: make(map[string]json.RawMessage)}
}

func NewParametersByName() *Parameters {
	return newParameters(ByName)
}

func NewParametersByPosition() *Parameters {
	return newParameters(ByPosition)
}

func (params *Parameters) GetType() Types {
	return params.paramType
}

func (params *Parameters) Set(key string, value json.RawMessage) (err error) {
	if params.paramType == ByPosition {
		_, err = strconv.Atoi(key)
		if err != nil {
			return err
		}
	}

	params.values[key] = value
	return nil
}

func (params *Parameters) SetString(key, value string) error {
	return params.Set(key, []byte(`"`+strings.ReplaceAll(value, "\"", "\\\"")+`"`))
}

func (params *Parameters) SetFloat(key string, value float64) error {
	return params.Set(key, []byte(strconv.FormatFloat(value, 'f', -1, 64)))
}

func (params *Parameters) SetInt(key string, value int) error {
	return params.Set(key, []byte(strconv.Itoa(value)))
}

func (params *Parameters) Get(key string) json.RawMessage {
	return params.values[key]
}

func (params *Parameters) GetString(key string) (value string, err error) {
	err = json.Unmarshal(params.Get(key), &value)
	return
}

func (params *Parameters) GetFloat(key string) (value float64, err error) {
	err = json.Unmarshal(params.Get(key), &value)
	return
}

func (params *Parameters) GetInt(key string) (value int, err error) {
	err = json.Unmarshal(params.Get(key), &value)
	return
}

func (params *Parameters) Serialize() (data json.RawMessage, err error) {
	if params.paramType == ByName {
		data, err = json.Marshal(params.values)
	} else {
		posMap := make(map[int]json.RawMessage)
		largestIndex := 0
		for key, value := range params.values {
			pos, err := strconv.Atoi(key)
			if err != nil {
				break
			}
			if pos > largestIndex {
				largestIndex = pos
			}
			posMap[pos] = value
		}

		rawData := make([]json.RawMessage, largestIndex+1)
		for key, value := range posMap {
			rawData[key] = value
		}

		data, err = json.Marshal(rawData)
	}

	if err != nil {
		return nil, err
	}
	return data, nil
}

func (params *Parameters) MarshalJSON() ([]byte, error) {
	return params.Serialize()
}

func (params *Parameters) UnmarshalJSON(jsonData []byte) (err error) {
	switch jsonData[0] {
	case '[':
		var data []json.RawMessage
		if err = json.Unmarshal(jsonData, &data); err != nil {
			break
		}

		for index, value := range data {
			params.values[strconv.Itoa(index)] = value
		}
		params.paramType = ByPosition
		break
	case '{':
		var data map[string]json.RawMessage
		if err = json.Unmarshal(jsonData, &data); err != nil {
			break
		}

		params.values = data
		params.paramType = ByName
		break
	default:
		err = errors.New("unable to parse parameters")
	}

	return
}
