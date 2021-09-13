package validators

import "strings"

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	errs := make([]string, len(v))
	for _, e := range v {
		errs = append(errs, e.Field+": "+e.Err.Error())
	}
	return strings.Join(errs, "\n")
}
