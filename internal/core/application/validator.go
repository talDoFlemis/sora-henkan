package application

import (
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
	once     sync.Once
)

// GetValidator returns a singleton validator instance
func GetValidator() *validator.Validate {
	once.Do(func() {
		validate = validator.New(validator.WithRequiredStructEnabled())
	})
	return validate
}

// ValidateStruct validates a struct and returns validation errors
func ValidateStruct(s interface{}) error {
	return GetValidator().Struct(s)
}
