package fileutils

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mholt/archiver"
)

// https://github.com/mholt/archiver

// Decompress decompresses files and choses the right filetype by its extension.
// It returns the decompressed filename.
func Decompress(src string) (string, error) {
	ext := filepath.Ext(src)
	if ext == "" {
		return "", fmt.Errorf("could not determine compression type for %s", src)
	}

	switch ext {
	case ".z", ".Z":
		return uncompressZ(src)
	default:
		return decompressFil(src)
	}
}

// decompress file using archiver modul
func decompressFil(src string) (string, error) {
	ext := filepath.Ext(src)
	dst := strings.TrimSuffix(src, ext)
	err := archiver.DecompressFile(src, dst)
	if err != nil {
		return "", fmt.Errorf("Could not uncompress %s: %v", src, err)
	}

	return dst, nil
}

// Uncompress .Z files made with "compress" tool.
// This cone be done with "gzip".
// The original file will be deleted
func uncompressZ(src string) (string, error) {
	tool, err := exec.LookPath("gzip")
	if err != nil {
		return "", err
	}

	ext := filepath.Ext(src)
	dst := strings.TrimSuffix(src, ext)

	cmd := exec.Command(tool, "-df", src)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("cmd %s failed: %v: %s", tool, err, stderr.Bytes())
	}

	if _, err := os.Stat(dst); os.IsNotExist(err) {
		return "", fmt.Errorf("%s failed: destination file %s does not exist: %v", tool, dst, err)
	}

	return dst, nil
}

// GzipDecompressor is an implementation of Decompressor that can
// decompress gzip files.
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
