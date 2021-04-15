package authmock

import "github.com/inspr/inspr/pkg/auth/models"

// MockAuth is the structure to mock the auth interface
type MockAuth struct {
	Err error
}

// NewMockAuth returns a new mock of the Auth interface
func NewMockAuth(err error) *MockAuth {
	return &MockAuth{Err: err}
}

// Validade - mock of the validate function
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

// Tokenize - mock of the tokenize function
func (ma *MockAuth) Tokenize(load models.Payload) ([]byte, error) {
	if ma.Err != nil {
		return []byte{}, ma.Err
	}
	return []byte("mock"), nil
}

// Refresh - mock of the refresh function
func (ma *MockAuth) Refresh(token []byte) ([]byte, error) {
	if ma.Err != nil {
		return []byte{}, ma.Err
	}
	return []byte("mock"), nil
}
