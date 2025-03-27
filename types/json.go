package types

import (
	"fmt"
)

type JsonError struct {
	Error string `json:"error"`
}

type JsonMessage struct {
	Message string `json:"message"`
}

type JsonMessageId struct {
	Message string `json:"message"`
	ID      int64  `json:"id"`
}

type JsonMessageFeed struct {
	Message string      `json:"message"`
	Feed    interface{} `json:"feed"`
}

type JsonMessageEntry struct {
	Message string      `json:"message"`
	Entry   interface{} `json:"entry"`
}

func NewJsonError(err error) JsonError {
	return JsonError{Error: err.Error()}
}

func NewMessage(message string) JsonMessage {
	return JsonMessage{Message: message}
}

func NewMessageF(format string, args ...any) JsonMessage {
	return JsonMessage{Message: fmt.Sprintf(format, args...)}
}
