package common

import "github.com/go-playground/validator/v10"

var validate = validator.New()

func ValidateRequest(request interface{}) error {
	err := validate.Struct(request)
	if err != nil {
		return err
	}
	return nil
}