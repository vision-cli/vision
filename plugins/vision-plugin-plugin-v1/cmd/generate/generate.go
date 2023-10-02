package generate

import (
	"embed"
	"encoding/json"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		newDir := "clone"
		err := os.MkdirAll(newDir, os.ModePerm)
		if err != nil {
			return err
		}
		err = fs.WalkDir(templateFiles, "template", func(path string, d fs.DirEntry, err error) error {
			// skip the top level template dir
			if path == "template" {
				return nil
			}
			switch {
			case path == "template":
				return nil
			case d.IsDir():
				// clone the dir
				return nil
			case filepath.Ext(path) == ".tmpl":
				// compute new path
			default:
				// clone
			}

			parts := strings.SplitN(path, string(os.PathSeparator), 2)
			if len(parts) != 2 {
				return fmt.Errorf("incorrect number of parts in path, expected 2 but got %d", len(parts))
			}
			newPath := parts[1] // path we want the new cloned file to be at

			// is this a templ file? if so drop the tmpl suffix
			var t *template.Template
			if filepath.Ext(newPath) == ".tmpl" {
				newPath = strings.TrimSuffix(newPath, filepath.Ext(newPath))
				t, err = template.ParseFiles(path)
				if err != nil {
					return err
				}
			}
			// strip template from all paths
			if d.IsDir() {
				return os.MkdirAll(filepath.Join(newDir, newPath), os.ModePerm)
			}
			fsrc, err := templateFiles.Open(path)
			if err != nil {
				return err
			}
			defer fsrc.Close()
			fdst, err := os.OpenFile(filepath.Join(newDir, newPath), os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
			if err != nil {
				return err
			}
			defer fdst.Close()
			if t != nil {
				err = t.Execute(fdst, struct{ ModuleName string }{ModuleName: "github.com/atos-digital/sample-plugin"})
				if err != nil {
					return err
				}
			} else {
				_, err = io.Copy(fdst, fsrc)
				if err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
		return json.NewEncoder(os.Stdout).Encode(map[string]any{
			"success": true,
		})
	},
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
