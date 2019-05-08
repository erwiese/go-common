package fileutils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Exists returns true if the file or directory exists, otherwise false.
func Exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true

	//info, err := os.Stat(filename)
	//return err == nil && !info.IsDir()
}

// CopyFile copies a file. If destination is a dir, the original filename will be kept.
// See https://opensource.com/article/18/6/copying-files-go
func CopyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	// if dest is a dir, use the src's filename
	if destFileStat, err := os.Stat(dst); !os.IsNotExist(err) {
		if destFileStat.Mode().IsDir() {
			_, srcFileName := filepath.Split(src)
			dst = filepath.Join(dst, srcFileName)
		}
	}

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

// RenameFile will rename the source to target using os function.
func RenameFile(oldpath, newpath string) error {
	return os.Rename(oldpath, newpath)
}

// MD5sum returns the computed MD5 checksum for the given file.
func MD5sum(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("Could not open file: %v", err)
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("Could not compute md5sum: %v", err)
	}

	hs := hex.EncodeToString(h.Sum(nil))

	return hs, nil
}

// RemoveAllContent removes all files and subdirs from a dir.
func RemoveAllContent(dir string) error {
	if len(dir) <= 3 {
		return fmt.Errorf("dir name is too short, do not clean: %s", dir)
	}
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()

	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}

	return nil
}
