package fileutils // github.com/erwiese/go-common/fileutils

import (
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mholt/archiver/v3"
)

// see https://github.com/klauspost/compress/tree/master/zstd#zstd

// FileCompressor can compress and decompress single files. As mholt/archiver does not support Z-compressed files,
// we make our own FileCompressor.
type FileCompressor struct {
	Src string // Source file, required.

	// Destination Dst can be empty, an existing directory or the explicit destination filename.
	// If it is empty or a directory, the destination filename will be built automatically.
	Dst string

	// Whether to overwrite existing files when creating files.
	OverwriteExisting bool

	// Whether the original file should be deleted after a successful de/compression.
	DeleteSource bool
}

// IsCompressed checks if a file is compressed, by its extension.
func IsCompressed(src string) bool {
	ext := filepath.Ext(src)
	if ext == "" {
		return false
	}

	if ext == ".z" || ext == ".Z" {
		return true
	}

	_, err := archiver.ByExtension(src)
	if err == nil {
		return true
	}

	return false
}

// Decompress the file and return the filename of the decompressed file. It choses the right filetype by its extension.
func (fc *FileCompressor) Decompress() (string, error) {
	if fc.Src == "" {
		return "", fmt.Errorf("Source file ist not defined")
	}
	ext := filepath.Ext(fc.Src)
	if ext == "" {
		return "", fmt.Errorf("could not determine compression type for %s", fc.Src)
	}

	if fc.Dst != "" {
		fi, err := os.Stat(fc.Dst)
		if err == nil { // otherwise dst is probably a file and does not exist
			if fi.IsDir() {
				fc.Dst = filepath.Join(fc.Dst, strings.TrimSuffix(filepath.Base(fc.Src), ext))
			}
		}
	} else {
		fc.Dst = strings.TrimSuffix(fc.Src, ext)
	}

	// Destination exists?
	if fi, err := os.Stat(fc.Dst); !os.IsNotExist(err) {
		if !fi.IsDir() && !fc.OverwriteExisting {
			fmt.Printf("Uncompress: destination file %s already exists. Exit\n", fc.Dst)
			return fc.Dst, nil
		}
	}

	switch ext {
	case ".z", ".Z":
		return fc.uncompressZ()
	default:
		return fc.decompressFil()
	}
}

// decompress file using mholt/archiver package
func (fc *FileCompressor) decompressFil() (string, error) {
	cIface, err := archiver.ByExtension(fc.Src)
	if err != nil {
		return "", err
	}

	dc, ok := cIface.(archiver.Decompressor)
	if !ok {
		return "", fmt.Errorf("format specified by source filename is not a recognized compression algorithm: %s", fc.Src)
	}
	comp := archiver.FileCompressor{Decompressor: dc, OverwriteExisting: fc.OverwriteExisting}
	err = comp.DecompressFile(fc.Src, fc.Dst)
	if err != nil {
		return "", fmt.Errorf("decompress %s: %v", fc.Src, err)
	}

	if _, err := os.Stat(fc.Dst); os.IsNotExist(err) {
		return "", fmt.Errorf("decompress %s: destination file does not exist", fc.Src)
	}

	if fc.DeleteSource {
		os.Remove(fc.Src)
	}

	return fc.Dst, nil
}

// Uncompress .Z files made with "compress" tool.
// This can be done with "gzip".
func (fc *FileCompressor) uncompressZ() (string, error) {
	tool, err := exec.LookPath("gzip")
	if err != nil {
		return "", err
	}

	options := "-dc"
	cmd := exec.Command(tool, options, fc.Src)
	err = RunCmdWithOutput(cmd, fc.Dst)
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(fc.Dst); os.IsNotExist(err) {
		return "", fmt.Errorf("%s failed: destination file %s does not exist: %v", tool, fc.Dst, err)
	}

	if fc.DeleteSource {
		os.Remove(fc.Src)
	}

	return fc.Dst, nil
}

// GzipDecompressor is an implementation of Decompressor that can decompress gzip files.
// From https://github.com/hashicorp/go-getter/blob/master/decompress_gzip.go
func gunzipAlternativ(src string) (string, error) {
	// Directory isn't supported at all
	// if dir {
	// 	return fmt.Errorf("gzip-compressed files can only unarchive to a single file")
	// }

	// If we're going into a directory we should make that first
	// if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
	// 	return "", err
	// }

	log.Printf("unzip file %s", src)

	dst := strings.TrimSuffix(src, ".gz")

	// File first
	f, err := os.Open(src)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// gzip compression is second
	gzipR, err := gzip.NewReader(f)
	if err != nil {
		return "", err
	}
	defer gzipR.Close()

	// Copy it out
	dstF, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer dstF.Close()

	_, err = io.Copy(dstF, gzipR)
	return dst, err
}

// http://blog.ralch.com/tutorial/golang-working-with-zip/
// https://socketloop.com/tutorials/unzip-compress-file-in-go
func unzipAlternativ(src string) (string, error) {
	reader, err := zip.OpenReader(src)
	if err != nil {
		return "", err
	}

	log.Printf("unzip file %s", src)

	dstPath := strings.TrimSuffix(src, filepath.Ext(src))

	// if err := os.MkdirAll(dstPath, 0755); err != nil {
	// 	return "", err
	// }

	dstFOrDir := ""
	for _, file := range reader.File {
		path := filepath.Join(dstPath, file.Name)
		if file.FileInfo().IsDir() {

			os.MkdirAll(path, file.Mode())
			if dstFOrDir == "" {
				dstFOrDir = path
			}
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return "", err
		}
		defer fileReader.Close()

		dstF, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return "", err
		}
		defer dstF.Close()

		if _, err := io.Copy(dstF, fileReader); err != nil {
			return "", err
		}
		if dstFOrDir == "" {
			dstFOrDir = path
		}
	}

	return dstFOrDir, nil
}

// ArchiveFiles creates an archive for a list of source files.
// The archive format is determined by the destinations' file extension, e.g. zip, tar.
func ArchiveFiles(sources []string, dest string) error {
	return archiver.Archive(sources, dest)
}

// CompressGzip compresses a file using gzip format and returns the compressed filename.
// Existing files will be overwritten.
// use: archiver.CompressFile(path, path+".gz") directly
/* func CompressGzip(path string) (string, error) {
	pathgz := path + ".gz"

	// try with gzip first, if installed
	if gzip, err := exec.LookPath("gzip"); err == nil {
		cmd := exec.Command(gzip, "-f", path) // -n
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		err = cmd.Run()
		if err != nil {
			return "", fmt.Errorf("gzip failed: %v: %s", err, stderr.Bytes())
		}
		if _, err := os.Stat(pathgz); os.IsNotExist(err) {
			return "", fmt.Errorf("gzip failed: %s: %s", "comressed file does not exist", pathgz)
		}
		return pathgz, nil
	}

	r, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer r.Close()

	// writer
	out, err := os.Create(pathgz)
	if err != nil {
		return "", err
	}
	defer out.Close()

	zw := gzip.NewWriter(out)
	_, err = io.Copy(zw, r)
	if err := zw.Close(); err != nil {
		return "", err
	}

	// check filesize > 0
	if finfo, err := os.Stat(pathgz); !os.IsNotExist(err) {
		if finfo.Size() < 1 {
			return "", fmt.Errorf("compressed file is empty: %s", pathgz)
		}
	}
	return pathgz, nil
} */
