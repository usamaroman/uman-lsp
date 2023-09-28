package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"lsp/messages"
	"lsp/mux"
)

const name = "uman lsp"

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdin, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	logger.Debug(fmt.Sprintf("%s started", name))

	m := mux.New(os.Stdin, os.Stdout, logger)
	m.HandleMethod("initialize", func(params json.RawMessage) (result any, err error) {
		var initializeParams messages.InitializeParams
		if err = json.Unmarshal(params, &initializeParams); err != nil {
			return result, err
		}

		result = messages.InitializeResult{
			Capabilities: messages.ServerCapabilities{
				TextDocumentSync: messages.TextDocumentSyncKindFull,
			},
			ServerInfo: &messages.ServerInfo{
				Name: name,
			},
		}

		return result, err
	})
	m.Log.Debug("registered method", slog.String("method", "initialize"))

	m.HandleNotification("initialized", func(params json.RawMessage) (err error) {
		return m.Notify("window/showMessage", messages.ShowMessageParams{
			Type:    messages.MessageTypeLog,
			Message: "added Uman support",
		})
	})
	m.Log.Debug("registered notification", slog.String("notification", "initialized"))

	for {
		if err := m.Process(); err != nil {
			return
		}
	}
}
