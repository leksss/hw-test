package validators

import (
	"errors"
)

var ErrMinValidatorShouldMore = errors.New("should be greater")

type MinValidator struct {
	Min int
}

func (v MinValidator) Validate(name string, val interface{}) *ValidationError {
	switch value := val.(type) {
	case int:
		if value < v.Min {
			return &ValidationError{name, ErrMinValidatorShouldMore}
		}
	case []int:
		for _, s := range value {
			if s < v.Min {
				return &ValidationError{name, ErrMinValidatorShouldMore}
			}
		}
	}
	return nil
}
