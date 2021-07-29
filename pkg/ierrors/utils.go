package ierrors

import (
	"errors"
)

// TODO add better description
// HasCode Compares error with error code
func HasCode(err error, code ErrCode) bool {
	e, ok := err.(*ierror)
	if !ok {
		return false
	}
	return (code & e.code) > 0
}

// TODO add better description
// Is Compares errors
func Is(source, target error) bool {
	// ierrors uses the func (ie *ierror) Is(err error) bool
	// to avaliate the two errors if they are an ierror type.
	// If not they will execute the standard comparison.
	return errors.Is(source, target)
}
