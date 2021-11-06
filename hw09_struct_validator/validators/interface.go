package validators

type Validator interface {
	Validate(string, interface{}) *ValidationError
}
