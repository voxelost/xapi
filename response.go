package xapi

import "github.com/google/uuid"

type response[T any] struct {
	Status bool `json:"status"`
	successResponse[T]
	errorResponse
	StreamSessionID *string   `json:"streamSessionId,omitempty"`
	MessageID       uuid.UUID `json:"customTag,omitempty"`
}

type successResponse[T any] struct {
	ReturnData *T `json:"returnData,omitempty"`
}

type errorResponse struct {
	ErrorCode        *string `json:"errorCode,omitempty"`
	ErrorDescription *string `json:"errorDescr,omitempty"`
}
