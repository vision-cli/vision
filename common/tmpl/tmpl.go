package tmpl

import (
	"errors"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/vision-cli/vision/common/file"
)

const templ_extension = ".tmpl"

// GenerateFS writes files to targetDir, mirroring the file system of templateFiles.
// Files with extension ".tmpl" will be templated with placeholder values parsed.
// SkipExisting will preserve the contents of files already in the targetDir.
func GenerateFS(templateFiles fs.FS, templateDir string, targetDir string, p any, skipExisting bool, t TmplWriter) error {
	return fs.WalkDir(templateFiles, templateDir, func(path string, d fs.DirEntry, err error) error {
		if errors.Is(err, fs.ErrPermission) {
			return fs.SkipDir
		}

		filename := strings.Replace(path, templateDir, targetDir, 1)
		filename = strings.Replace(filename, templ_extension, "", 1)

		if d.IsDir() {
			return t.CreateDir(filename)
		}

		if skipExisting && file.Exists(filename) {
			return nil
		}

		if IsTemplate(path) {
			return t.WriteTemplatedFS(path, filename, templateFiles, p)
		}
		return t.WriteExactFS(path, filename, templateFiles)
	})
}

func IsTemplate(path string) bool {
	return filepath.Ext(path) == templ_extension
}
