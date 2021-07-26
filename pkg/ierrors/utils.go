package ierrors

import "errors"

// TODO add better description
// Is Compares errors
func Is(source, target error) bool {

	// check if is another type of error inside the error stack
	if errors.Is(source, target) {
		return true
	}

	// converts target to ierror structure, if possible
	t, ok := target.(*ierror)
	if !ok {
		return false
	}

	// converts source to ierror structure, if possible
	s, ok := target.(*ierror)
	if !ok {
		return false
	}

	return s.code&t.code > 0
}

// TODO add better description
// HasCode Compares error with error code
func HasCode(err error, code ErrCode) bool {
	e, ok := err.(*ierror)
	if !ok {
		return false
	}
	return (code & e.code) > 0
}
