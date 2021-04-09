package ierrors

func HasCode(target error, code InsprErrorCode) bool {
	t, ok := target.(*InsprError)
	if !ok {
		return false
	}
	return t.Code&code > 0
}

func IsIerror(target error) bool {
	t, ok := target.(*InsprError)
	if !ok {
		return false
	}

	// going through all the errors of the ierror pkg
	if t.Code > 0 {
		return true
	}
	return false
}
