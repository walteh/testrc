// Code generated by mockery v2.33.2. DO NOT EDIT.

package mockery

import mock "github.com/stretchr/testify/mock"

// MockDummy_tests is an autogenerated mock type for the Dummy type
type MockDummy_tests struct {
	mock.Mock
}

type MockDummy_tests_Expecter struct {
	mock *mock.Mock
}

func (_m *MockDummy_tests) EXPECT() *MockDummy_tests_Expecter {
	return &MockDummy_tests_Expecter{mock: &_m.Mock}
}

// NewMockDummy_tests creates a new instance of MockDummy_tests. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockDummy_tests(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockDummy_tests {
	mock := &MockDummy_tests{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
