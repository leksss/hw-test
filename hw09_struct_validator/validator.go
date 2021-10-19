package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/leksss/hw-test/hw09_struct_validator/validators"
)

const validatorTagName = "validate"

const (
	validatorNameLen    = "len"
	validatorNameMin    = "min"
	validatorNameMax    = "max"
	validatorNameRegexp = "regexp"
	validatorNameIn     = "in"
)

var ErrExpectStruct = errors.New("expected a struct")

func Validate(v interface{}) (error, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Struct {
		return nil, ErrExpectStruct
	}

	errs := validators.ValidationErrors{}

	t := rv.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(validatorTagName)
		if tag == "" || tag == "-" {
			continue
		}

		value := rv.Field(i).Interface()
		validatorStrings := strings.Split(tag, "|")
		if len(validatorStrings) == 0 {
			continue
		}

		for _, vs := range validatorStrings {
			validator, err := createValidator(vs)
			if err != nil {
				return nil, err
			}

			if err := validator.Validate(field.Name, value); err != nil {
				errs = append(errs, *err)
			}
		}
	}

	return errs, nil
}

func createValidator(validatorString string) (validators.Validator, error) {
	validatorName, validatorValue := parseValidatorString(validatorString)

	switch validatorName {
	case validatorNameLen:
		validator := validators.LenValidator{}
		length, err := strconv.Atoi(validatorValue)
		if err != nil {
			return nil, fmt.Errorf("length validator: %w", err)
		}
		validator.Length = length
		return validator, nil

	case validatorNameIn:
		validator := validators.InValidator{}
		in := strings.Split(validatorValue, ",")
		if len(in) == 0 {
			return nil, fmt.Errorf("invalid in validator")
		}
		validator.In = in
		return validator, nil

	case validatorNameMin:
		validator := validators.MinValidator{}
		min, err := strconv.Atoi(validatorValue)
		if err != nil {
			return nil, fmt.Errorf("min validator: %w", err)
		}
		validator.Min = min
		return validator, nil

	case validatorNameMax:
		validator := validators.MaxValidator{}
		max, err := strconv.Atoi(validatorValue)
		if err != nil {
			return nil, fmt.Errorf("invalid max validator: %w", err)
		}
		validator.Max = max
		return validator, nil

	case validatorNameRegexp:
		validator := validators.RegexpValidator{}
		if validatorValue == "" {
			return nil, fmt.Errorf("invalid regexp validator")
		}
		validator.Regexp = validatorValue
		return validator, nil
	}
	return validators.DefaultValidator{}, nil
}

func parseValidatorString(validatorString string) (string, string) {
	validatorData := strings.Split(validatorString, ":")
	if len(validatorData) != 2 {
		return "", ""
	}
	return validatorData[0], validatorData[1]
}
