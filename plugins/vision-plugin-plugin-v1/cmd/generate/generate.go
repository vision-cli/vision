package generate

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

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
	newDir := "clone"
	err := os.MkdirAll(newDir, os.ModePerm)
	if err != nil {
		return err
	}
	return fs.WalkDir(templateFiles, "template", func(path string, d fs.DirEntry, err error) error {
		// skip the top level template dir
		newPath := strings.TrimPrefix(path, "template/")
		switch {
		case path == "template":
			return nil
		case d.IsDir():
			// if it is a dir then create it
			return os.MkdirAll(filepath.Join("clone", newPath), os.ModePerm)
		case filepath.Ext(newPath) == ".tmpl":
			// if this is a template file then remove the .tmpl suffix
			trimmedNewPath := strings.TrimSuffix(newPath, filepath.Ext(newPath))
			// create the file
			fsrc, err := templateFiles.Open(path)
			if err != nil {
				return fmt.Errorf("opening from templateFiles: %w", err)
			}
			defer fsrc.Close()
			fdst, err := os.OpenFile(filepath.Join("clone", trimmedNewPath), os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
			if err != nil {
				return fmt.Errorf("[tmpl] opening from clone: %v, %w", filepath.Join("clone", trimmedNewPath), err)
			}
			defer fdst.Close()
			_, err = io.Copy(fdst, fsrc)
			return err
		default:
			// clone
			// create the file
			fsrc, err := templateFiles.Open(path)
			if err != nil {
				return fmt.Errorf("opening from templateFiles: %w", err)
			}
			defer fsrc.Close()
			fdst, err := os.OpenFile(filepath.Join("clone", newPath), os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
			if err != nil {
				return fmt.Errorf("[clone] opening from clone: %w", err)
			}
			defer fdst.Close()
			_, err = io.Copy(fdst, fsrc)
			return err
		}
	})
}

func cloneDir(dst, src string) error {
	return nil
}

func cloneFile(dst, src string) error {
	return nil
}

// clone and execute a template
func cloneExecTmpl() error {
	return nil
}
