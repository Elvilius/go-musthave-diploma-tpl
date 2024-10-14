package mocks

type MockToken struct{}

func NewMockToken() *MockToken {
	return &MockToken{}
}

func (t *MockToken) GenerateTokenForUser(userID int) (string, error) {
	return "secret", nil
}
