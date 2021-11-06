package validators

import (
	"errors"
	"regexp"
)

var ErrRegexpValidatorNotMatch = errors.New("value should match regexp")

type RegexpValidator struct {
	Regexp string
}

func (v RegexpValidator) Validate(name string, val interface{}) *ValidationError {
	regExp := regexp.MustCompile(v.Regexp)
	if !regExp.MatchString(val.(string)) {
		return &ValidationError{name, ErrRegexpValidatorNotMatch}
	}
	return nil
}
