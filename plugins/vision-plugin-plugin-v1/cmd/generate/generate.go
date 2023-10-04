package generate

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

//go:embed all:template
var templateFiles embed.FS

var GenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "the plugin version",
	Long:  "ditto",
	RunE:  run,
}

func run(cmd *cobra.Command, args []string) error {
	err := cloneDir("clone")
	if err != nil {
		return fmt.Errorf("cloning directory: %w", err)
	}

	return fs.WalkDir(templateFiles, "template", func(path string, d fs.DirEntry, err error) error {
		newPath := strings.TrimPrefix(path, "template/")

		switch {
		case path == "template": // skip the top level template dir
			return nil
		case d.IsDir(): // if it is a dir then create it
			return cloneDir(filepath.Join("clone", newPath))
		case filepath.Ext(newPath) == ".tmpl":
			err := cloneTmplFile(newPath, path)
			if err != nil {
				return fmt.Errorf("cloning template files: %w", err)
			}
			return nil
		default:
			cloneFile(newPath, path)
			if err != nil {
				return fmt.Errorf("cloning files: %w", err)
			}
			return nil
		}
	})
}

func cloneDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func cloneFile(dst, src string) error {
	fsrc, err := templateFiles.Open(src)
	if err != nil {
		return fmt.Errorf("opening from templateFiles: %w", err)
	}
	defer fsrc.Close()
	fdst, err := os.OpenFile(filepath.Join("clone", dst), os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return fmt.Errorf("[clone] opening from clone: %w", err)
	}
	defer fdst.Close()
	_, err = io.Copy(fdst, fsrc)
	return err
}

func cloneTmplFile(dst, src string) error {
	// if this is a template file then remove the .tmpl suffix
	trimmedNewPath := strings.TrimSuffix(dst, filepath.Ext(dst))
	// create the file
	return cloneFile(trimmedNewPath, src)
}

type ReadmeFile struct {
	PluginName string
}

func cloneExecTmpl(fileName string, f fs.File) error {

	b, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	tmplEx, err := template.New(fileName).Parse(string(b))
	if err != nil {
		return err
	}

	p1 := ReadmeFile{
		PluginName: "ExamplePlugin",
	}

	return tmplEx.Execute(os.Stdout, p1)
}
