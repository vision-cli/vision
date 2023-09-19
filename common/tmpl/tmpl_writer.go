package tmpl

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
	"text/template"

	"github.com/vision-cli/vision/common/cases"
	"github.com/vision-cli/vision/common/file"
)

type TmplWriter interface {
	WriteTemplatedFS(templatePath string, targetPath string, templateFiles fs.FS, p interface{}) error
	WriteExactFS(templatePath string, targetPath string, templateFiles fs.FS) error
	CreateDir(path string) error
}

type OsTmplWriter struct{}

func NewOsTmpWriter() TmplWriter {
	return OsTmplWriter{}
}

func (OsTmplWriter) WriteTemplatedFS(templatePath string, targetPath string, templateFiles fs.FS, p interface{}) error {
	t, err := newTemplateFS(templatePath, templateFiles)
	if err != nil {
		return fmt.Errorf("creating template for %s: %w", targetPath, err)
	}

	newF, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer newF.Close()

	err = giveExecPermissionsIfScript(newF)
	if err != nil {
		return err
	}

	return t.Execute(newF, p)
}

func (OsTmplWriter) WriteExactFS(templatePath string, targetPath string, templateFiles fs.FS) error {
	src, err := templateFiles.Open(templatePath)
	if err != nil {
		return err
	}
	defer src.Close()

	return copyFile(src, targetPath)
}

func (OsTmplWriter) CreateDir(path string) error {
	return file.CreateDir(path)
}

// Returns a template with the standard function map
func New(name string, text string) (*template.Template, error) {
	funcs := template.FuncMap{
		"Pascal": cases.Pascal,
		"Camel":  cases.Camel,
		"Snake":  cases.Snake,
		"Kebab":  cases.Kebab,
	}

	return template.New(name).Funcs(funcs).Parse(text)
}

func TmplToString(text string, tokens interface{}) (string, error) {
	tmpl, err := New("temp", text)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, tokens)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// newTemplateFS returns a template, with all the templating functions, from the path.
func newTemplateFS(path string, fsys fs.FS) (*template.Template, error) {
	f, err := fsys.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening template file: %w", err)
	}
	defer f.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, f)
	if err != nil {
		return nil, fmt.Errorf("copying bytes from template file: %w", err)
	}

	t, err := New(path, buf.String())
	if err != nil {
		return nil, fmt.Errorf("creating template from file: %w", err)
	}

	return t, nil
}

func copyFile(src fs.File, targetPath string) error {
	dst, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	err = giveExecPermissionsIfScript(dst)
	if err != nil {
		return err
	}

	_, err = io.Copy(dst, src)
	return err
}

// Changes mode for files with ".sh" extension to allow execution
func giveExecPermissionsIfScript(f *os.File) error {
	if strings.HasSuffix(f.Name(), ".sh") {
		return file.GiveExecPermissions(f)
	}
	return nil
}
