package common

var ErrBadRequestShouldBind = "Invalid request body"
type ResponseError struct {
	Message string `json:"error"`
}

func NewResponseError(message string) *ResponseError {
	return &ResponseError{
		Message: message,
	}
}