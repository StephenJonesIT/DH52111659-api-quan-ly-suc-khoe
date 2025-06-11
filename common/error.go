package common

var ErrBadRequestShouldBind = "Invalid request body"
var ErrRecordNotFound = "Record not found"
var ErrInvalidPhoneNumber = "The phone number is invalid"
var ErrExpertNotFound = "Expert not found"
var ErrProgramNotFound = "Program not found"
var ErrInternalServerError = "Internal server error"
var ErrUnauthorized = "Unauthorized"

type ResponseError struct {
	Message string `json:"error"`
}

func NewResponseError(message string) *ResponseError {
	return &ResponseError{
		Message: message,
	}
}