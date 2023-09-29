package messages

const DidChangeTextDocumentNotification = "textDocument/didChange"

type DidChangeTextDocumentParams struct {
	TextDocument   VersionedTextDocumentIdentifier  `json:"textDocument"`
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

type VersionedTextDocumentIdentifier struct {
	Version int    `json:"version"`
	URI     string `json:"uri"`
}

type TextDocumentContentChangeEvent struct {
	Range *Range `json:"range"`
	Text  string `json:"text"`
}

func NewPosition(line uint, ch rune) Position {
	return Position{
		Line:      line,
		Character: ch,
	}
}

type Position struct {
	Line      uint `json:"line"`
	Character rune `json:"character"`
}

type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

type Location struct {
	URI   string `json:"uri"`
	Range Range  `json:"range"`
}
