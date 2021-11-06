package validators

import (
	"errors"
)

var ErrMaxValidatorShouldBeLess = errors.New("should be less")

type MaxValidator struct {
	Max int
}

func (v MaxValidator) Validate(name string, val interface{}) *ValidationError {
	switch value := val.(type) {
	case int:
		if value > v.Max {
			return &ValidationError{name, ErrMaxValidatorShouldBeLess}
		}
	case []int:
		for _, s := range value {
			if s > v.Max {
				return &ValidationError{name, ErrMaxValidatorShouldBeLess}
			}
		}
	}
	return nil
}
