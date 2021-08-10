package group

import "github.com/gelleson/packup/pkg/validators"

type CreateGroupInput struct {
	Name string `validate:"required"`
}

func (input CreateGroupInput) Validate() error {
	return validators.Struct(input)
}
