package domains

import "github.com/favclip/golidator"

var validator = golidator.NewValidator()

func Validate(obj interface{}) error {
	return validator.Validate(obj)
}
