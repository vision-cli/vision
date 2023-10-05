package generate

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
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

// structs need to go into some sort of struct place?? API - is that what it's called??
type visionJson struct {
	PluginName string
	PluginData samplePlugin `json:"helloworld"`
}

type samplePlugin struct {
	ModuleName string `json:"module_name"`
	Key2       []int  `json:"key2"`
}

func run(cmd *cobra.Command, args []string) error {
	vj, err := openVisionJson()
	if err != nil {
		return fmt.Errorf("opening vision.json: %w", err)
	}

	err = cloneDir("clone")
	if err != nil {
		return fmt.Errorf("cloning directory: %w", err)
	}

	return fs.WalkDir(templateFiles, "template", func(path string, d fs.DirEntry, err error) error {
		newPath := filepath.Join("clone", strings.TrimPrefix(path, "template/"))

		switch {
		case path == "template": // skip the top level template dir
			return nil
		case d.IsDir(): // if it is a dir then create it
			return cloneDir(newPath)
		case filepath.Ext(newPath) == ".tmpl":
			err := cloneExecTmpl(path, newPath, vj)
			if err != nil {
				return fmt.Errorf("cloning template files: %w", err)
			}

			return nil
		default:
			cloneFile(path, newPath)
			if err != nil {
				return fmt.Errorf("cloning files: %w", err)
			}
			return nil
		}
	})
}

func openVisionJson() (*visionJson, error) {
	// TODO(luke): create a "VISIONPATH" env variable and look for that?
	f, err := os.OpenFile("../../vision.json", os.O_RDWR, 0444)
	if err != nil {
		return nil, fmt.Errorf("opening config file: %w", err)
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("reading bytes: %w", err)
	}

	// set default PluginName value
	jsonData := visionJson{
		PluginName: "helloworld",
	}
	if err = json.Unmarshal(b, &jsonData); err != nil {
		return nil, fmt.Errorf("unmarshalling json: %w", err)
	}

	return &jsonData, nil
}

func cloneDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func cloneFile(src, dst string) error {
	fsrc, err := templateFiles.Open(src)
	if err != nil {
		return fmt.Errorf("opening from templateFiles: %w", err)
	}
	defer fsrc.Close()
	fdst, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return fmt.Errorf("[clone] opening from clone: %w", err)
	}
	defer fdst.Close()
	_, err = io.Copy(fdst, fsrc)
	return err
}

func cloneExecTmpl(src, dst string, vj *visionJson) error {
	// open file and read it
	trimmedNewPath := strings.TrimSuffix(dst, filepath.Ext(dst))
	err := cloneFile(src, trimmedNewPath)
	if err != nil {
		return fmt.Errorf("cloning file: %w", err)
	}
	f, err := os.OpenFile(trimmedNewPath, os.O_RDWR, 0644) // only enable reading mode as we do not need to write anything
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}
	err = f.Truncate(0)
	if err != nil {
		return fmt.Errorf("truncating: %w", err)
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("seeking: %w", err)
	}

	tmplEx, err := template.New("templateFile").Parse(string(b))
	if err != nil {
		return fmt.Errorf("creating template file: %w", err)
	}

	fmt.Printf("vision json: %+v\n", vj)

	return tmplEx.Execute(f, vj)
	// return nil
}
