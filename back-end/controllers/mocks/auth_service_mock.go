package mocks

import "github.com/stretchr/testify/mock"

type AuthServiceMock struct {
	mock.Mock
}

func (m *AuthServiceMock) Authenticate(email string, password string) (*uint, error) {
	args := m.Called(email, password)
	return args.Get(0).(*uint), args.Error(1)
}
