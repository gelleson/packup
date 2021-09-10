package dto

import "github.com/gelleson/packup/pkg/validators"

type CreateUserInput struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8"`
	GroupId  uint
}

func (input CreateUserInput) Validate() error {
	return validators.Struct(input)
}

func (input CreateUserInput) HasGroup() bool {
	return input.GroupId != 0
}
