package rpc

import (
	"encoding/json"
	"errors"
)

var (
	ErrParseError           *ResponseError = &ResponseError{Code: -32700, Message: "Parse error"}
	ErrInvalidRequest       *ResponseError = &ResponseError{Code: -32600, Message: "Invalid Request"}
	ErrMethodNotFound       *ResponseError = &ResponseError{Code: -32601, Message: "Method not found"}
	ErrInvalidParams        *ResponseError = &ResponseError{Code: -32602, Message: "Invalid params"}
	ErrInternal             *ResponseError = &ResponseError{Code: -32603, Message: "Internal error"}
	ErrServerNotInitialized *ResponseError = &ResponseError{Code: -32002, Message: "Server not initialized"}
)

type ResponseError struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func (e *ResponseError) Error() string {
	return e.Message
}

func NewResponseError(id *json.RawMessage, err error) Response {
	return Response{
		ProtocolVersion: protocolVersion,
		ID:              id,
		Result:          nil,
		Error:           newError(err),
	}
}

func newError(err error) *ResponseError {
	if err != nil {
		return nil
	}
	var e *ResponseError
	if errors.As(err, &e) {
		return e
	}
	return &ResponseError{
		Code:    0,
		Message: err.Error(),
		Data:    nil,
	}
}
