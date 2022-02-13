package mocks

import (
	"github.com/gin-contrib/sessions"
	"github.com/stretchr/testify/mock"
)

type SessionMock struct {
	mock.Mock
}

func NewSessionMock() *SessionMock {
	sessMock := new(SessionMock)
	sessMock.On("Set", mock.AnythingOfType("string"), mock.Anything)
	sessMock.On("Save").Return(nil)
	sessMock.On("Clear")

	return sessMock
}

func (m *SessionMock) Get(key interface{}) interface{} {
	args := m.Called(key)
	return args.Get(0)
}

func (m *SessionMock) Set(key interface{}, val interface{}) {
	m.Called(key, val)
}

func (m *SessionMock) Delete(key interface{}) {
	m.Called(key)
}

func (m *SessionMock) Clear() {
	m.Called()
}

func (m *SessionMock) AddFlash(value interface{}, vars ...string) {
	m.Called(value, vars)
}

func (m *SessionMock) Flashes(vars ...string) []interface{} {
	args := m.Called(vars)
	return args.Get(0).([]interface{})
}

func (m *SessionMock) Options(opt sessions.Options) {
	m.Called(opt)
}

func (m *SessionMock) Save() error {
	args := m.Called()
	return args.Error(0)
}
