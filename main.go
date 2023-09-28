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
	logFile, err := os.Create(".umanlsp.log")
	if err != nil {
		slog.Error("failed to create log output file", slog.Any("error", err))
		os.Exit(1)
	}
	defer logFile.Close()

	logger := slog.New(slog.NewJSONHandler(logFile, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	defer func() {
		if r := recover(); r != nil {
			logger.Error("panic", slog.Any("recovered", r))
		}
	}()

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
