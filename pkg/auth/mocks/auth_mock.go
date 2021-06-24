package authmock

import (
	"inspr.dev/inspr/pkg/auth"
)

// MockAuth is the structure to mock the auth interface
type MockAuth struct {
	Err error
}

// NewMockAuth returns a new mock of the Auth interface
func NewMockAuth(err error) *MockAuth {
	return &MockAuth{Err: err}
}

// Validate - mock of the validate function
func (ma *MockAuth) Validate(token []byte) (*auth.Payload, []byte, error) {
	if ma.Err != nil {
		return nil, []byte{}, ma.Err
	}
	return &auth.Payload{
		UID: "uid",
		Permissions: map[string][]string{
			"scope_1": {auth.CreateChannel},
			"scope_2": {},
		},
		Refresh:    []byte("refresh"),
		RefreshURL: "refresh_url",
	}, []byte("mock"), nil
}

// Init - mock of the tokenize function
func (ma *MockAuth) Init(s string, load auth.Payload) ([]byte, error) {
	if ma.Err != nil {
		return []byte{}, ma.Err
	}
	return []byte("mock"), nil
}

// Tokenize - mock of the tokenize function
func (ma *MockAuth) Tokenize(load auth.Payload) ([]byte, error) {
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
