package phd2_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/goastro/phd2"
)

func TestLoop(t *testing.T) {
	type testCase struct {
		name string

		writeReturnVal []interface{}
		readReturnVal  []interface{}

		expectedResult bool
		expectedErr    error
	}

	testCases := []testCase{
		testCase{
			name: "Good",
			writeReturnVal: []interface{}{
				1,
				nil,
			},
			readReturnVal: []interface{}{
				[]byte{0},
				nil,
			},
			expectedResult: true,
			expectedErr:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := &MockDialer{}
			conn := &MockConn{}

			d.On("Dial", "tcp", "127.0.0.1:4300").Return(conn, nil)
			conn.On("Write", []byte{0x13}).Return(tc.writeReturnVal...)
			conn.On("Read", make([]byte, 1)).Return(tc.readReturnVal...)

			c := phd2.NewSocketClient(d)

			err := c.Connect("127.0.0.1", 4300)
			require.NoError(t, err)

			success, err := c.Loop()

			assert.Equal(t, tc.expectedResult, success)

			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			d.AssertExpectations(t)
			conn.AssertExpectations(t)
		})
	}
}
