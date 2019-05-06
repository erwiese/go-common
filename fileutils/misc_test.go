package fileutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMD5sum(t *testing.T) {
	sum, err := MD5sum("testdata/ctab340k.18d.Z")
	assert.NoError(t, err)
	assert.Equal(t, "e32566b9227a25216e44658646968707", sum, "md5sum")
}
