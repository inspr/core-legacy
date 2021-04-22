package request

type mockAuth struct {
	err error
}

func (ma mockAuth) GetToken() ([]byte, error) {
	if ma.err != nil {
		return []byte{}, ma.err
	}
	return []byte("bearer mock_token"), nil
}

func (ma mockAuth) SetToken([]byte) error {
	if ma.err != nil {
		return ma.err
	}
	return nil
}
