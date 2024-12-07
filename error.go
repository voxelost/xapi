package xapi

type ApiError struct {
	Code    string
	Message string
}

func (e ApiError) Error() string {
	return e.Message
}
