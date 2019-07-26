package fileutils

import (
	"os"
	"testing"
)

func init() {
	_ = os.Mkdir("testdata/tmp/", 0700) // os.ModePerm (0777)
	RemoveAllContent("testdata/tmp/")
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

func TestDecompress(t *testing.T) {
	// Copy source files tmp dir first
	zFil := "testdata/ctab340k.18d.Z"
	_, err := CopyFile(zFil, "testdata/tmp")
	if err != nil {
		t.Errorf("could not copy Z-file to tmp dir: %v", err)
	}
	zFil = "testdata/tmp/ctab340k.18d.Z"

	gzFil := "testdata/DENT00BEL_R_20183401000_01H_30S_MO.crx.gz"
	_, err = CopyFile(gzFil, "testdata/tmp")
	if err != nil {
		t.Errorf("could not copy gz-file to tmp dir: %v", err)
	}
	gzFil = "testdata/tmp/DENT00BEL_R_20183401000_01H_30S_MO.crx.gz"

	os.Mkdir("testdata/tmp/2", 0700)

	type fields struct {
		Src               string
		Dst               string
		OverwriteExisting bool
		DeleteSource      bool
	}

	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "T10 decompress *.Z  simple",
			fields: fields{
				Src:               zFil,
				Dst:               "",
				OverwriteExisting: false,
				DeleteSource:      false,
			},
			want:    "testdata/tmp/ctab340k.18d",
			wantErr: false,
		},
		{
			name: "T11 decompress *.gz simple",
			fields: fields{
				Src:               gzFil,
				Dst:               "",
				OverwriteExisting: false,
				DeleteSource:      false,
			},
			want:    "testdata/tmp/DENT00BEL_R_20183401000_01H_30S_MO.crx",
			wantErr: false,
		},
		{
			name: "T20 decompress *.Z with destination which is a dir",
			fields: fields{
				Src:               zFil,
				Dst:               "testdata/tmp/2",
				OverwriteExisting: false,
				DeleteSource:      false,
			},
			want:    "testdata/tmp/2/ctab340k.18d",
			wantErr: false,
		},
		{
			name: "T21 decompress *.gz with destination which is a dir",
			fields: fields{
				Src:               gzFil,
				Dst:               "testdata/tmp/2",
				OverwriteExisting: false,
				DeleteSource:      false,
			},
			want:    "testdata/tmp/2/DENT00BEL_R_20183401000_01H_30S_MO.crx",
			wantErr: false,
		},
		{
			name: "T30 decompress *.Z do not overwrite existing file",
			fields: fields{
				Src:               zFil,
				Dst:               "testdata/tmp/2/ctab340k.18d",
				OverwriteExisting: false,
				DeleteSource:      false,
			},
			want:    "testdata/tmp/2/ctab340k.18d",
			wantErr: false,
		},
		{
			name: "T31 decompress *.gz do not overwrite existing file",
			fields: fields{
				Src:               gzFil,
				Dst:               "testdata/tmp/2/DENT00BEL_R_20183401000_01H_30S_MO.crx",
				OverwriteExisting: false,
				DeleteSource:      false,
			},
			want:    "testdata/tmp/2/DENT00BEL_R_20183401000_01H_30S_MO.crx",
			wantErr: false,
		},
		{
			name: "T31 decompress *.gz overwrite existing file",
			fields: fields{
				Src:               gzFil,
				Dst:               "testdata/tmp/2/DENT00BEL_R_20183401000_01H_30S_MO.crx",
				OverwriteExisting: true,
				DeleteSource:      false,
			},
			want:    "testdata/tmp/2/DENT00BEL_R_20183401000_01H_30S_MO.crx",
			wantErr: false,
		},
		{
			name: "T40 decompress *.Z with destination which is a file",
			fields: fields{
				Src:               zFil,
				Dst:               "testdata/tmp/wursti",
				OverwriteExisting: false,
				DeleteSource:      false,
			},
			want:    "testdata/tmp/wursti",
			wantErr: false,
		},
		{
			name: "T41 decompress *.gz with destination which is a file",
			fields: fields{
				Src:               gzFil,
				Dst:               "testdata/tmp/bratli",
				OverwriteExisting: false,
				DeleteSource:      false,
			},
			want:    "testdata/tmp/bratli",
			wantErr: false,
		},

		// last test: delete source
		{
			name: "T50 decompress *.Z delete source",
			fields: fields{
				Src:               zFil,
				Dst:               "testdata/tmp/wursti",
				OverwriteExisting: true,
				DeleteSource:      true,
			},
			want:    "testdata/tmp/wursti",
			wantErr: false,
		},
		{
			name: "T51 decompress *.gz delete source",
			fields: fields{
				Src:               gzFil,
				Dst:               "testdata/tmp/bratli",
				OverwriteExisting: true,
				DeleteSource:      true,
			},
			want:    "testdata/tmp/bratli",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := &FileCompressor{
				Src:               tt.fields.Src,
				Dst:               tt.fields.Dst,
				OverwriteExisting: tt.fields.OverwriteExisting,
				DeleteSource:      tt.fields.DeleteSource,
			}
			got, err := fc.Decompress()
			if (err != nil) != tt.wantErr {
				t.Errorf("FileCompressor.Decompress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FileCompressor.Decompress() = %v, want %v", got, tt.want)
			}
			if tt.fields.DeleteSource {
				if _, err := os.Stat(tt.fields.Src); !os.IsNotExist(err) {
					t.Errorf("FileCompressor.Decompress() error: source was not deleted")
				}
			}
		})
	}
}
