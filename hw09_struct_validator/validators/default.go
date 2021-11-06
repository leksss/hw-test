package validators

type DefaultValidator struct{}

func (v DefaultValidator) Validate(name string, val interface{}) *ValidationError {
	return nil
}
