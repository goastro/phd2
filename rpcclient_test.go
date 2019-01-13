package phd2_test

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/goastro/phd2"
)

func TestRPCClient(t *testing.T) {
	c := phd2.NewRPCClient(&net.Dialer{})

	err := c.Connect("127.0.0.1", 4400)

	assert.NoError(t, err)

	time.Sleep(5 * time.Second)
	t.Fail()
}
