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
	for i := 0; i < 32; i++ {
		checkErr := t.Code & (1 << i)
		if checkErr > 0 {
			return true
		}
	}
	return false
}
