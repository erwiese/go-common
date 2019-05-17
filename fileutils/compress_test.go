package fileutils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	_ = os.Mkdir("testdata/tmp/", 0700) // os.ModePerm (0777)
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

func TestArchiveFiles(t *testing.T) {
	type args struct {
		files   []string
		zipfile string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "t1",
			args: args{
				files:   []string{"testdata/ctab340k.18d.Z", "testdata/DENT00BEL_R_20183401000_01H_30S_MO.crx.gz"},
				zipfile: "testdata/tmp/rinex-t1.zip",
			},
			wantErr: false,
		},
		{
			name: "t2",
			args: args{
				files:   []string{"testdata/ctab340k.18d.Z", "testdata/DENT00BEL_R_20183401000_01H_30S_MO.crx.gz"},
				zipfile: "testdata/tmp/rinex-t1.tar",
			},
			wantErr: false,
		},
		{
			name: "t3-NotExistingInputFile",
			args: args{
				files:   []string{"testdata/ctab340k.18d", "testdata/DENT00BEL_R_20183401000_01H_30S_MO.crx.gz"},
				zipfile: "testdata/tmp/rinex-t2.zip",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ArchiveFiles(tt.args.files, tt.args.zipfile); (err != nil) != tt.wantErr {
				t.Errorf("ArchiveFiles() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
