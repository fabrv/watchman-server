package utils

import "github.com/go-playground/validator/v10"

type ErrorResponse struct {
	FailedField string `json:"failed_field"`
	Tag         string `json:"tag"`
	Value       string `json:"value"`
}

var validate = validator.New()

func ValidateStruct(s interface{}) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, &ErrorResponse{
				FailedField: err.Field(),
				Tag:         err.Tag(),
				Value:       err.Param(),
			})
		}
	}
	return errors
}
