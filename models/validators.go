package models

import (
	"github.com/go-playground/validator"
)

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func SetDefaultValue(s interface{}) {
	switch model := s.(type) {
	case *CreateUserRequest:
		if model.UserType == "" {
			model.UserType = GuestUser
		}
	case *CreateBoardRequest:
		if model.Template == "" {
			model.Template = Agile
		}
	}
}

func ValidateStruct(s interface{}) []*ErrorResponse {
	var errors []*ErrorResponse

	SetDefaultValue(s)

	validation := validator.New()
	err := validation.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
