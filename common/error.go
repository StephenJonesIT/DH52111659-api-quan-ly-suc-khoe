package common

type ResponseError struct {
	Message string `json:"error"`
}

func NewResponseError(message string) *ResponseError {
	return &ResponseError{
		Message: message,
	}
}