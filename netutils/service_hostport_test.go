package netutils

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHostPortLookup(t *testing.T) {
	assert := require.New(t)

	publicip, err := FindPublicIPv4()
	assert.NoError(err)

	assert.Equal("127.0.0.1:1234", ServiceHostPort(&net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 1234,
	}))
	assert.Equal(fmt.Sprintf("%s:1234", publicip.String()), ServiceHostPort(&net.TCPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: 1234,
	}))

	assert.Equal("127.0.0.1:1234", ServiceHostPort(&net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 1234,
	}))
	assert.Equal(fmt.Sprintf("%s:1234", publicip.String()), ServiceHostPort(&net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: 1234,
	}))

}
