package toolbox

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringByteSize(t *testing.T) {
	assert := require.New(t)

	assert.Equal("64 byte", StringByteSize(64))
	assert.Equal("512 byte", StringByteSize(512))
	assert.Equal("1.00 MiB", StringByteSize(1024*1024))
	assert.Equal("1.25 MiB", StringByteSize(1024*1024+1024*256))
	assert.Equal("1.00 GiB", StringByteSize(1024*1024*1024))
	assert.Equal("1.00 TiB", StringByteSize(1024*1024*1024*1024))
	assert.Equal("1.00 PiB", StringByteSize(1024*1024*1024*1024*1024))
	assert.Equal("1.00 EiB", StringByteSize(1024*1024*1024*1024*1024*1024))
	assert.Equal("4.00 EiB", StringByteSize(1024*1024*1024*1024*1024*1024*4))
}
