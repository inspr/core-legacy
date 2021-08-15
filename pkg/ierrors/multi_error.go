package ierrors

// MultiError is a type that handles multiple errors that can happen in a process
type MultiError struct {
	Errors []error `yaml:"errors" json:"slices"`
	Code   ErrCode `yaml:"code"   json:"code"`
}

// Error return the errors in the multierror concatenated with new lines
func (e *MultiError) Error() (ret string) {
	if e.Empty() {
		return ""
	}
	for _, err := range e.Errors[:len(e.Errors)-1] {
		ret += err.Error() + "\n"
	}
	ret += e.Errors[len(e.Errors)-1].Error()
	return
}

// Add adds an error to the multi error
func (e *MultiError) Add(err error) {
	if err != nil {
		e.Errors = append(e.Errors, err)
		if ierr, ok := err.(*ierror); ok {
			e.Code |= ierr.code
		}
	}
}

// Empty returns if there is no error in the multierror
func (e *MultiError) Empty() bool {
	return len(e.Errors) <= 0
}
