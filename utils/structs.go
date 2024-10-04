package utils

import "fmt"

type JsonError struct {
	Error string `json:"error"`
}

type JsonMessage struct {
	Message string `json:"message"`
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
