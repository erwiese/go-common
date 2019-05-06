package fileutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	RemoveAllContent("testdata/tmp/")
}

func TestDecompress(t *testing.T) {
	assert := assert.New(t)

	// Gzip file
	// Copy to tmp dir first
	filePath := "testdata/DENT00BEL_R_20183401000_01H_30S_MO.crx.gz"
	_, err := CopyFile(filePath, "testdata/tmp")
	if err != nil {
		t.Fatalf("Could not copy file to tmp dir: %v", err)
	}
	filePath = "testdata/tmp/DENT00BEL_R_20183401000_01H_30S_MO.crx.gz"

	dst, err := Decompress(filePath)
	if err != nil {
		t.Fatalf("Could not decompress file: %v", err)
	}
	assert.Equal("testdata/tmp/DENT00BEL_R_20183401000_01H_30S_MO.crx", dst, "decompress gzip")
}

func TestDecompressZ(t *testing.T) {
	assert := assert.New(t)

	// zip file
	// Copy to tmp dir first
	filePath := "testdata/ctab340k.18d.Z"
	_, err := CopyFile(filePath, "testdata/tmp")
	if err != nil {
		t.Fatalf("Could not copy file to tmp dir: %v", err)
	}
	filePath = "testdata/tmp/ctab340k.18d.Z"

	dst, err := Decompress(filePath)
	if err != nil {
		t.Fatalf("Could not decompress file: %v", err)
	}
	assert.Equal("testdata/tmp/ctab340k.18d", dst, "decompress *.Z")
}
