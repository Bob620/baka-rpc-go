package errors

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewParseError() *RPCError {
	return &RPCError{
		Code:    -32700,
		Message: "Parse error",
	}
}

func NewInvalidRequest() *RPCError {
	return &RPCError{
		Code:    -32600,
		Message: "Invalid Request",
	}
}

func NewMethodNotFound() *RPCError {
	return &RPCError{
		Code:    -32601,
		Message: "Method not found",
	}
}

func NewInvalidParams() *RPCError {
	return &RPCError{
		Code:    -32602,
		Message: "Invalid params",
	}
}

func NewInternalError() *RPCError {
	return &RPCError{
		Code:    -32603,
		Message: "Internal error",
	}
}

func NewGenericError(err string) *RPCError {
	return &RPCError{
		Code:    -32000,
		Message: err,
	}
}
