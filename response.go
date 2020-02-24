package main

import "encoding/json"

type SuccessResponse struct {
	JsonRpc string          `json:"jsonrpc"`
	Id      *string         `json:"id"`
	Result  json.RawMessage `json:"result"`
}

type ErrorResponse struct {
	JsonRpc string    `json:"jsonrpc"`
	Id      *string   `json:"id"`
	Error   *RPCError `json:"error"`
}

func NewSuccessResponse(id *string, result json.RawMessage) json.RawMessage {
	message, _ := json.Marshal(SuccessResponse{
		JsonRpc: "2.0",
		Result:  result,
		Id:      id,
	})

	return message
}

func NewErrorResponse(id *string, error *RPCError) json.RawMessage {
	message, _ := json.Marshal(ErrorResponse{
		JsonRpc: "2.0",
		Error:   error,
		Id:      id,
	})

	return message
}
