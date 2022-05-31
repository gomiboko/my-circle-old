package mocks

import "github.com/stretchr/testify/mock"

type SessionServiceMock struct {
	mock.Mock
}

func (m *SessionServiceMock) Authenticate(email string, password string) (*uint, error) {
	args := m.Called(email, password)
	return args.Get(0).(*uint), args.Error(1)
}
