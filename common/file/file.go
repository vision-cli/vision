package file

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

var Osstat = os.Stat
var Osremove = os.Remove
var Osreaddir = os.ReadDir
var Osopen = os.Open
var Osremoveall = os.RemoveAll
var Osgetwd = os.Getwd
var Osgetenv = os.Getenv

// CreateDir creates a directory, along with any necessary parents.
// If path is already a file that is not a directory,
// CreateDir will remove the file and create a directory in its place.
func CreateDir(path string) error {
	info, err := os.Stat(path)
	if err == nil && !info.IsDir() {
		os.Remove(path)
	}
	return os.MkdirAll(path, os.ModePerm)
}

// DeleteIfEmptyDir deletes path if it is an accessible empty directory.
func DeleteIfEmptyDir(path string) {
	if isAcessibleEmptyDir(path) {
		_ = Osremove(path)
	}
}

// DeleteEmptyDirs Deletes any immediate child directories that are empty.
func DeleteEmptyDirs(targetDir string) error {
	fds, err := Osreaddir(targetDir)
	if err != nil {
		return err
	}
	for _, fd := range fds {
		path := filepath.Join(targetDir, fd.Name())
		DeleteIfEmptyDir(path)
	}
	return nil
}

func isAcessibleEmptyDir(path string) bool {
	dir, err := Osopen(path)
	if err != nil {
		return false
	}
	defer dir.Close()

	_, err = dir.Readdirnames(1)
	return errors.Is(err, io.EOF)
}

// Exists returns true if path exists and is accessible
func Exists(path string) bool {
	_, err := Osstat(path)
	// other errors or nil imply existence (e.g. ErrPermission)
	return !(errors.Is(err, os.ErrNotExist) || errors.Is(err, os.ErrInvalid))
}

// RemoveNamed removes each named file relative to dir. Non-existent paths are ignored.
func RemoveNamed(dir string, named ...string) error {
	var err error

	for _, n := range named {
		err = Osremoveall(filepath.Join(dir, n))
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}

	return nil
}

// changes file permissions to executable
func GiveExecPermissions(f *os.File) error {
	return f.Chmod(0755) //nolint:gomnd //permission mask has inherent meaning (rwx r-x r-x)
}

// Wrapper around os.Getwd
func GetWorkingDir() (string, error) {
	return Osgetwd()
}

// Wrapper around os.ReadDir
func ReadDir(path string) ([]os.DirEntry, error) {
	return Osreaddir(path)
}

// Wrapper around os.GetEnv
func GetEnv(key string) string {
	return Osgetenv(key)
}

// Wrapper around os.Open
func Open(name string) (*os.File, error) {
	return Osopen(name)
}
