package ali_mns

import fasthttp "github.com/valyala/fasthttp"
import mock "github.com/stretchr/testify/mock"

// MockMNSClient is an autogenerated mock type for the MNSClient type
type mockMNSClient struct {
	mock.Mock
}

// getAccountID provides a mock function with given fields:
func (_m *mockMNSClient) getAccountID() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// getRegion provides a mock function with given fields:
func (_m *mockMNSClient) getRegion() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Send provides a mock function with given fields: method, headers, message, resource, opts
func (_m *mockMNSClient) Send(method Method, headers map[string]string, message interface{}, resource string, opts ...Option) (*fasthttp.Response, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, method, headers, message, resource)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *fasthttp.Response
	if rf, ok := ret.Get(0).(func(Method, map[string]string, interface{}, string, ...Option) *fasthttp.Response); ok {
		r0 = rf(method, headers, message, resource, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*fasthttp.Response)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(Method, map[string]string, interface{}, string, ...Option) error); ok {
		r1 = rf(method, headers, message, resource, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetProxy provides a mock function with given fields: url
func (_m *mockMNSClient) SetProxy(url string) {
	_m.Called(url)
}
