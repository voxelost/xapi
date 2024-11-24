package xapi

import "github.com/google/uuid"

type command[T any] struct {
	Command   string    `json:"command"`
	Arguments T         `json:"arguments,omitempty"`
	MessageID uuid.UUID `json:"customTag,omitempty"`
}

func newCommand[T any](cmd string, data T) *command[T] {
	return &command[T]{
		Command:   cmd,
		Arguments: data,
		MessageID: uuid.New(),
	}
}
