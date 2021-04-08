package authMock

import "gitlab.inspr.dev/inspr/core/pkg/auth/models"

type MockAuth struct {
	Err error
}

func NewMockAuth(err error) *MockAuth {
	return &MockAuth{Err: err}
}

func (ma *MockAuth) Validade(token []byte) (models.Payload, []byte, error) {
	if ma.Err != nil {
		return models.Payload{}, []byte{}, ma.Err
	}
	return models.Payload{
		UID:        "uid",
		Role:       0,
		Scope:      []string{"scope_1", "scope_2"},
		Refresh:    "refresh",
		RefreshURL: "refresh_url",
	}, []byte("mock"), nil
}
func (ma *MockAuth) Tokenize(load models.Payload) ([]byte, error) {
	if ma.Err != nil {
		return []byte{}, ma.Err
	}
	return []byte("mock"), nil
}
func (ma *MockAuth) Refresh(token []byte) ([]byte, error) {
	if ma.Err != nil {
		return []byte{}, ma.Err
	}
	return []byte("mock"), nil
}
