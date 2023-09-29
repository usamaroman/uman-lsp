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

	fileURIToContents := make(map[string]string)
	documentUpdates := make(chan messages.TextDocumentItem, 10)
	go func() {
		for doc := range documentUpdates {
			fileURIToContents[doc.URI] = doc.Text
		}
	}()

	m.HandleMethod(messages.InitializeMethod, func(params json.RawMessage) (result any, err error) {
		var initializeParams messages.InitializeParams
		if err = json.Unmarshal(params, &initializeParams); err != nil {
			logError(m, err)
			return result, err
		}

		result = messages.InitializeResult{
			Capabilities: messages.ServerCapabilities{
				TextDocumentSync: messages.TextDocumentSyncKindFull,
				CompletionProvider: &messages.CompletionOpts{
					TriggerCharacters: []string{"с", "в", ":", "ч", "д", "л", "и"},
				},
			},
			ServerInfo: &messages.ServerInfo{
				Name: name,
			},
		}

		return
	})
	m.Log.Debug("registered method", slog.String("method", messages.InitializeMethod))

	m.HandleNotification(messages.InitializedNotification, func(params json.RawMessage) (err error) {
		return m.Notify("window/showMessage", messages.ShowMessageParams{
			Type:    messages.MessageTypeInfo,
			Message: "added Uman support",
		})
	})
	m.Log.Debug("registered notification", slog.String("notification", messages.InitializedNotification))

	m.HandleNotification(messages.DidOpenTextDocumentNotification, func(rawParams json.RawMessage) (err error) {
		var params messages.DidOpenTextDocumentParams

		if err = json.Unmarshal(rawParams, &params); err != nil {
			logError(m, err)
			return
		}

		documentUpdates <- params.TextDocument
		return nil
	})
	m.Log.Debug("registered notification", slog.String("notification", messages.DidOpenTextDocumentNotification))

	m.HandleNotification(messages.DidChangeTextDocumentNotification, func(rawParams json.RawMessage) (err error) {
		var params messages.DidChangeTextDocumentParams

		if err = json.Unmarshal(rawParams, &params); err != nil {
			logError(m, err)
			return
		}

		documentUpdates <- messages.TextDocumentItem{
			URI:     params.TextDocument.URI,
			Version: params.TextDocument.Version,
			Text:    params.ContentChanges[0].Text,
		}
		return nil
	})
	m.Log.Debug("registered notification", slog.String("notification", messages.DidChangeTextDocumentNotification))

	m.HandleMethod(messages.CompletionRequestMethod, func(params json.RawMessage) (result any, err error) {
		var completionParams messages.CompletionParams
		if err = json.Unmarshal(params, &completionParams); err != nil {
			logError(m, err)
			return result, err
		}

		var r = []messages.CompletionItem{
			{
				Label:         "число",
				Kind:          messages.CompletionItemKindKeyword,
				Detail:        "числовой тип данных",
				Documentation: "числовой тип данных",
			},
			{
				Label:         "строка",
				Kind:          messages.CompletionItemKindKeyword,
				Detail:        "строковой тип данных",
				Documentation: "строковой тип данных",
			},
			{
				Label:         "истина",
				Kind:          messages.CompletionItemKindKeyword,
				Detail:        "булевое значение",
				Documentation: "булевое значение",
			},
			{
				Label:         "ложь",
				Kind:          messages.CompletionItemKindKeyword,
				Detail:        "булевое значение",
				Documentation: "булевое значение",
			},
			{
				Label:         "создать",
				Kind:          messages.CompletionItemKindKeyword,
				Detail:        "ключевое слово для создания переменных",
				Documentation: "ключевое слово для создания переменных",
			},
			{
				Label:         "вывести",
				Kind:          messages.CompletionItemKindFunction,
				Detail:        "вывести значение переменной",
				Documentation: "вывести значение переменной",
			},
			{
				Label:         "длина",
				Kind:          messages.CompletionItemKindFunction,
				Detail:        "вывести длину",
				Documentation: "вывести длину",
			},
			{
				Label:         "если",
				Kind:          messages.CompletionItemKindKeyword,
				Detail:        "условный оператор",
				Documentation: "условный оператор",
			},
			{
				Label:         "иначе",
				Kind:          messages.CompletionItemKindKeyword,
				Detail:        "условный оператор",
				Documentation: "условный оператор",
			},
			{
				Label:         "функция",
				Kind:          messages.CompletionItemKindKeyword,
				Detail:        "функция",
				Documentation: "функция",
			},
		}

		return r, err
	})
	m.Log.Debug("registered method", slog.String("method", messages.CompletionRequestMethod))

	for {
		if err := m.Process(); err != nil {
			logError(m, err)
			return
		}
	}
}

func logError(m *mux.Mux, err error) {
	m.Log.Error("got error", slog.String("error", err.Error()))
}
