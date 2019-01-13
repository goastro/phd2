package phd2_test

import (
	"net"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockDialer struct {
	mock.Mock
}

func (m *MockDialer) Dial(network, address string) (net.Conn, error) {
	args := m.Called(network, address)

	conn := args.Get(0)
	err := args.Error(1)

	if conn == nil {
		return nil, err
	}

	return conn.(net.Conn), err
}

type MockConn struct {
	mock.Mock
}

func (m *MockConn) Read(b []byte) (n int, err error) {
	args := m.Called(b)

	if bytesToRead, ok := args.Get(0).([]byte); ok {
		var i int
		for i = 0; i < len(bytesToRead) && i < len(b); i++ {
			b[i] = bytesToRead[i]
		}

		return i, args.Error(1)
	}

	return 0, args.Error(1)
}

func (m *MockConn) Write(b []byte) (n int, err error) {
	args := m.Called(b)
	return args.Int(0), args.Error(1)
}

func (m *MockConn) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockConn) LocalAddr() net.Addr {
	args := m.Called()
	addr := args.Get(0)

	if addr == nil {
		return nil
	}

	return addr.(net.Addr)
}

func (m *MockConn) RemoteAddr() net.Addr {
	args := m.Called()
	addr := args.Get(0)

	if addr == nil {
		return nil
	}

	return addr.(net.Addr)
}

func (m *MockConn) SetDeadline(t time.Time) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m *MockConn) SetReadDeadline(t time.Time) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m *MockConn) SetWriteDeadline(t time.Time) error {
	args := m.Called(t)
	return args.Error(0)
}
