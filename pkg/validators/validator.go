package validators

import "github.com/go-playground/validator/v10"

var v = validator.New()

func Struct(i interface{}) error {
	return v.Struct(i)
}
