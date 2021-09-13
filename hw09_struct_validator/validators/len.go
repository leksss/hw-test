package validators

import (
	"errors"
)

var ErrLenValidatorMustBeExact = errors.New("string must be exact chars long")

type LenValidator struct {
	Length int
}

func (v LenValidator) Validate(name string, val interface{}) *ValidationError {
	switch value := val.(type) {
	case string:
		if len(value) != v.Length {
			return &ValidationError{name, ErrLenValidatorMustBeExact}
		}
	case []string:
		for _, s := range value {
			if len(s) != v.Length {
				return &ValidationError{name, ErrLenValidatorMustBeExact}
			}
		}
	}
	return nil
}
