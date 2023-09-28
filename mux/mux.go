package mux

import (
	"bufio"
	"encoding/json"
	"log/slog"
	"os"
	"sync"

	"lsp/rpc"
)

type NotificationHandler func(params json.RawMessage) (err error)
type MethodHandler func(params json.RawMessage) (result any, err error)

type Mux struct {
	reader               *bufio.Reader
	writer               *bufio.Writer
	notificationHandlers map[string]NotificationHandler
	methodHandlers       map[string]MethodHandler
	writeLock            *sync.Mutex
	Log                  *slog.Logger
}

func New(r, w *os.File, log *slog.Logger) *Mux {
	return &Mux{
		reader:               bufio.NewReader(r),
		writer:               bufio.NewWriter(w),
		notificationHandlers: make(map[string]NotificationHandler),
		methodHandlers:       make(map[string]MethodHandler),
		writeLock:            &sync.Mutex{},
		Log:                  log,
	}
}

func (m *Mux) Process() (err error) {
	req, err := rpc.Read(m.reader)
	if err != nil {
		return err
	}
	m.Log.Info("got request", slog.Any("request", req))

	go func(req rpc.Request) {
		if req.IsNotification() {
			if nh, ok := m.notificationHandlers[req.Method]; ok {
				nh(req.Params)
			}
		} else {
			mh, ok := m.methodHandlers[req.Method]
			if !ok {
				m.write(rpc.NewResponseError(req.ID, rpc.ErrMethodNotFound))
				return
			}

			result, err := mh(req.Params)
			if err != nil {
				m.write(rpc.NewResponseError(req.ID, err))
				return
			}
			m.write(rpc.NewResponse(req.ID, result))
		}
	}(req)

	return
}

func (m *Mux) HandleMethod(name string, handler MethodHandler) {
	m.methodHandlers[name] = handler
}

func (m *Mux) HandleNotification(name string, handler NotificationHandler) {
	m.notificationHandlers[name] = handler
}

func (m *Mux) Notify(method string, params any) error {
	n := rpc.Notification{
		ProtocolVersion: "2.0",
		Method:          method,
		Params:          params,
	}

	return m.write(n)
}

func (m *Mux) write(msg rpc.Message) error {
	m.writeLock.Lock()
	defer m.writeLock.Unlock()

	m.Log.Info("sent response", slog.Any("response", msg))

	return rpc.Write(m.writer, msg)
}
