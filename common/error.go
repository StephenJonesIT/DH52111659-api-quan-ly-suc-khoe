package common

var ErrBadRequestShouldBind = "Invalid request body"
var ErrRecordNotFound = "Record not found"
var ErrInvalidPhoneNumber = "The phone number is invalid"
type ResponseError struct {
	Message string `json:"error"`
}

func NewResponseError(message string) *ResponseError {
	return &ResponseError{
		Message: message,
	}
}