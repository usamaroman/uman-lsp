package messages

type MessageType int

const (
	_ MessageType = iota
	MessageTypeError
	MessageTypeWarning
	MessageTypeInfo
	MessageTypeLog
	MessageTypeDebug
)

type ShowMessageParams struct {
	Type    MessageType `json:"type"`
	Message string      `json:"message"`
}
