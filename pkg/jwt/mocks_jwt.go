package jwt

type MockToken struct{}

func NewMockToken() *MockToken {
	return &MockToken{}
}

func (t *MockToken) GenerateTokenForUser(_ int) (string, error) {
	return "secret", nil
}
