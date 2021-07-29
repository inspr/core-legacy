package ierrors

// HasCode tries to convert the error interface to an ierror structure and then
// assert if the structure contains the ErrCode provided as an argument
func HasCode(err error, code ErrCode) bool {
	e, ok := err.(*ierror)
	if !ok {
		return false
	}
	return (code & e.code) > 0
}
