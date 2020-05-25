package parameters

import (
	"encoding/json"
	"errors"
	"strconv"
)

type Types string

const (
	ByName     Types = "byName"
	ByPosition Types = "byPosition"
)

type Parameters struct {
	paramType Types
	values    map[string]Param
}

func newParameters(paramType Types, params []Param) *Parameters {
	if params == nil {
		return &Parameters{paramType: paramType, values: make(map[string]Param)}
	}

	values := make(map[string]Param)
	for _, param := range params {
		values[param.GetName()] = param
	}

	return &Parameters{paramType: paramType, values: values}
}

func NewParametersByName(params []Param) *Parameters {
	return newParameters(ByName, params)
}

func NewParametersByPosition(params []Param) *Parameters {
	for index, param := range params {
		param.SetName(strconv.Itoa(index))
	}

	return newParameters(ByPosition, params)
}

func (params *Parameters) Length() int {
	return len(params.values)
}

func (params *Parameters) GetType() Types {
	return params.paramType
}

func (params *Parameters) Set(key string, value Param) (err error) {
	if params.paramType == ByPosition {
		_, err = strconv.Atoi(key)
		if err != nil {
			return err
		}
	}

	params.values[key] = value
	return nil
}

func (params *Parameters) Get(key string) Param {
	return params.values[key]
}

func (params *Parameters) Serialize() (data json.RawMessage, err error) {
	if params.paramType == ByName {
		data, err = json.Marshal(params.values)
	} else {
		posMap := make(map[int]Param)
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

		rawData := make([]Param, largestIndex+1)
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

		params.values = make(map[string]Param)
		for index, value := range data {
			pos := strconv.Itoa(index)
			params.values[pos] = &GenericParam{pos, value, false, value}
		}
		params.paramType = ByPosition
		break
	case '{':
		var data map[string]json.RawMessage
		if err = json.Unmarshal(jsonData, &data); err != nil {
			break
		}

		params.values = make(map[string]Param)
		for key, value := range data {
			params.values[key] = &GenericParam{key, value, false, value}
		}

		params.paramType = ByName
		break
	default:
		err = errors.New("unable to parse parameters")
	}

	return
}
