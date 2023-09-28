package rpc

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/textproto"
	"strconv"
)

const protocolVersion = "2.0"

var ErrInvalidContentLengthHeader = errors.New("missing or invalid Content-Length header")

type Message interface {
	IsJsonRPC() bool
}

type Request struct {
	ProtocolVersion string           `json:"jsonrpc"`
	ID              *json.RawMessage `json:"id"`
	Method          string           `json:"method"`
	Params          json.RawMessage  `json:"params"`
}

func (r Request) IsJsonRPC() bool {
	return r.ProtocolVersion == protocolVersion
}

func (r Request) IsNotification() bool {
	return r.ID == nil
}

type Response struct {
	ProtocolVersion string           `json:"jsonrpc"`
	ID              *json.RawMessage `json:"id"`
	Result          any              `json:"result"`
	Error           *ResponseError   `json:"error"`
}

func NewResponse(id *json.RawMessage, result any) Response {
	return Response{
		ProtocolVersion: protocolVersion,
		ID:              id,
		Result:          result,
		Error:           nil,
	}
}

func (r Response) IsJsonRPC() bool {
	return r.ProtocolVersion == protocolVersion
}

type Notification struct {
	ProtocolVersion string `json:"jsonrpc"`
	Method          string `json:"method"`
	Params          any    `json:"params"`
}

func (n Notification) IsJsonRPC() bool {
	return n.ProtocolVersion == protocolVersion
}

func Read(r *bufio.Reader) (req Request, err error) {
	header, err := textproto.NewReader(r).ReadMIMEHeader()
	if err != nil {
		return
	}

	contentLength, err := strconv.ParseInt(header.Get("Content-Length"), 10, 64)
	if err != nil {
		return req, ErrInvalidContentLengthHeader
	}

	err = json.NewDecoder(io.LimitReader(r, contentLength)).Decode(&req)
	if err != nil {
		return
	}

	if !req.IsJsonRPC() {
		return req, ErrInvalidRequest
	}

	return
}

func Write(w *bufio.Writer, msg Message) (err error) {
	body, err := json.Marshal(msg)
	if err != nil {
		return
	}

	headers := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(body))
	if _, err = w.WriteString(headers); err != nil {
		return
	}
	if _, err = w.Write(body); err != nil {
		return
	}

	return w.Flush()
}
