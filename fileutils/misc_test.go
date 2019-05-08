package fileutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExists(t *testing.T) {
	// dir
	ok := Exists("testdata")
	assert.True(t, ok)

	ok = Exists("testdataaaa")
	assert.False(t, ok)

	// file
	ok = Exists("testdata/ctab340k.18d.Z")
	assert.True(t, ok)

	ok = Exists("testdata/ctab340k.18d.Ziiiii")
	assert.False(t, ok)
}

func TestMD5sum(t *testing.T) {
	sum, err := MD5sum("testdata/ctab340k.18d.Z")
	assert.NoError(t, err)
	assert.Equal(t, "e32566b9227a25216e44658646968707", sum, "md5sum")
}
