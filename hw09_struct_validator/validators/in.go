package validators

import (
	"errors"
	"reflect"
	"strconv"
)

var ErrInValidatorShouldBeInList = errors.New("should be in list")

type InValidator struct {
	In []string
}

func (v InValidator) Validate(name string, val interface{}) *ValidationError {
	switch value := val.(type) {
	case int, []int:

		switch t := value.(type) {
		case int:
			value = []int{t}
		default:
		}

		for _, i := range value.([]int) {
			if !containsInt(i, v.convertToIntIn()) {
				return &ValidationError{name, ErrInValidatorShouldBeInList}
			}
		}

	default:
		rv := reflect.ValueOf(value)
		s := rv.String()
		if !containsString(s, v.In) {
			return &ValidationError{name, ErrInValidatorShouldBeInList}
		}
	}

	return nil
}

func (v InValidator) convertToIntIn() []int {
	intIn := make([]int, len(v.In))
	for _, v := range v.In {
		i, _ := strconv.Atoi(v)
		intIn = append(intIn, i)
	}
	return intIn
}

func containsInt(e int, s []int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func containsString(e string, s []string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
