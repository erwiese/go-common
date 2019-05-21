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

func TestIsSane(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "t1-illegal",
			path: "../../../somewhat",
			want: false,
		},
		{
			name: "t2-tooLong",
			path: "/home/user/somewhat/igssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssss",
			want: false,
		},
		{
			name: "t3",
			path: "archive/tmp/myFile.txt",
			want: true,
		},
		{
			name: "t4",
			path: "/home/user/tmp/duzipfe",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSane(tt.path); got != tt.want {
				t.Errorf("IsSane() = %v, want %v", got, tt.want)
			}
		})
	}
}
