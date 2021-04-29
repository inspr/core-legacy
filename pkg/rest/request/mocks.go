package request

type mockAuth struct {
	errGet error
	errSet error
}

func (ma mockAuth) GetToken() ([]byte, error) {
	if ma.errGet != nil {
		return []byte{}, ma.errGet
	}
	return []byte("bearer mock_token"), nil
}

func (ma mockAuth) SetToken([]byte) error {
	if ma.errSet != nil {
		return ma.errSet
	}
	return nil
}
