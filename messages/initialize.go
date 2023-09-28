package messages

type InitializeParams struct {
	// ClientInfo provides information about client
	ClientInfo *ClientInfo `json:"clientInfo"`

	// Capabilities provided by client
	Capabilities ClientCapabilities `json:"capabilities"`
}

type InitializeResult struct {
	// Capabilities provided by LSP server
	Capabilities ServerCapabilities `json:"capabilities"`

	// ServerInfo provides information about server
	ServerInfo *ServerInfo `json:"serverInfo"`
}

type ClientInfo struct {
	Name    string  `json:"name"`
	Version *string `json:"version"`
}

type ServerInfo struct {
	Name    string  `json:"name"`
	Version *string `json:"version"`
}

type ClientCapabilities struct{}

type ServerCapabilities struct {
	TextDocumentSync   TextDocumentSyncKind `json:"textDocumentSync"`
	CompletionProvider *CompletionOpts      `json:"completionProvider,omitempty"`
}

type TextDocumentSyncKind int

const (
	TextDocumentSyncKindNone TextDocumentSyncKind = iota
	TextDocumentSyncKindFull
	TextDocumentSyncKindIncremental
)

type CompletionOpts struct {
	TriggerCharacters []string `json:"triggerCharacters,omitempty"`
}
