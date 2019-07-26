package fileutils

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

// RunCmdWithOutput runs the specified command and writes its output to the given file.
func RunCmdWithOutput(cmd *exec.Cmd, outfile string) error {
	f, err := os.Create(outfile)
	if err != nil {
		return fmt.Errorf("Could not create file %s: %v", outfile, err)
	}
	defer f.Close()
	writer := bufio.NewWriter(f)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("RunCmd error: %v", err)
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err = cmd.Start() // as new OS process
	if err != nil {
		return fmt.Errorf("starting cmd failed: %v", err)
	}

	go io.Copy(writer, stdout)

	//see https://stackoverflow.com/questions/10385551/get-exit-code-go
	if err := cmd.Wait(); err != nil { // no timeout, see context
		return fmt.Errorf("cmd failed: %v: %s", err, stderr.Bytes())
	}
	writer.Flush()

	fmt.Printf("file %s created\n", outfile)

	return nil
}

// IsSane returns true if the filepath seems to be sane for further processing.
// Especially useful for checking form inputs.
func IsSane(path string) bool {
	if len(path) > 150 {
		fmt.Printf("path is too long: %s", path)
		return false
	}

	if strings.Contains(path, "..") {
		fmt.Printf("illegal characters: %s", path)
		return false
	}

	// TODO add checks

	return true
}
