package tmplutils

import (
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
)

const TmplExt = ".gohtml"

// WalkDir walks the file path and returns a slice of paths to files with the template extension (.gohtml).
func WalkDir(path string) ([]string, error) {
	paths := []string{}
	return paths, filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() || filepath.Ext(path) != TmplExt {
				return nil
			}
			paths = append(paths, path)
			return nil
		})
}

// ParseDir parses the files in the file path and returns a template.
func ParseDir(templateName, path string) (*template.Template, error) {
	tmpl := template.New(templateName)
	paths, err := WalkDir(path)
	if err != nil {
		return nil, err
	}
	return tmpl.ParseFiles(paths...)
}

// WalkFS walks the file system and returns a slice of paths to files with the template extension (.gohtml).
func WalkFS(fileSys fs.FS) ([]string, error) {
	paths := []string{}
	return paths, fs.WalkDir(fileSys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || filepath.Ext(path) != TmplExt {
			return nil
		}
		paths = append(paths, path)
		return nil
	})
}

// ParseFS parses the files in the file system and returns a template.
func ParseFS(templateName string, fileSys fs.FS) (*template.Template, error) {
	tmpl := template.New(templateName)
	paths, err := WalkFS(fileSys)
	if err != nil {
		return nil, err
	}
	return tmpl.ParseFS(fileSys, paths...)
}
