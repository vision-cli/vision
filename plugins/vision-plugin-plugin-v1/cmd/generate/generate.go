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

	"github.com/vision-cli/vision/plugins/vision-plugin-plugin-v1/cmd/initialise"
)

//go:embed all:template
var templateFiles embed.FS

// TODO (luke): improve description
var GenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate the plugins code",
	Long:  "generate the plugins code using the vision.json config",
	RunE:  generate,
}

type success struct {
	Success bool `json:"success"`
}

func generate(cmd *cobra.Command, args []string) error {
	err := run(cmd, args)
	jEnc := json.NewEncoder(os.Stdout)
	if err != nil {
		_ = jEnc.Encode(success{Success: false})
		return fmt.Errorf("generating template: %w", err)
	}

	err = jEnc.Encode(success{Success: true})
	if err != nil {
		return fmt.Errorf("encoding JSON response: %w", err)
	}
	return nil
}

func run(cmd *cobra.Command, args []string) error {
	var vPath string
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %w", err)
	}
	if len(args) < 1 {
		vPath = filepath.Join(wd, "vision.json")
	} else {
		vPath = args[0]
	}

	vj, err := openVisionJson(vPath)
	if err != nil {
		return fmt.Errorf("opening vision.json: %w", err)
	}

	pluginDir := strings.TrimSuffix(vPath, "vision.json")

	err = cloneDir(pluginDir)
	if err != nil {
		return fmt.Errorf("cloning directory: %w", err)
	}

	return fs.WalkDir(templateFiles, "template", func(path string, d fs.DirEntry, err error) error {
		newPath := filepath.Join(pluginDir, strings.TrimPrefix(path, "template/"))

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

func openVisionJson(vPath string) (*initialise.PluginConfig, error) {
	f, err := os.OpenFile(vPath, os.O_RDWR, 0444)
	if err != nil {
		return nil, fmt.Errorf("opening config file: %w", err)
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("reading bytes: %w", err)
	}

	var jsonData initialise.PluginConfig
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

func cloneExecTmpl(src, dst string, vj *initialise.PluginConfig) error {
	// open file and read it
	trimmedNewPath := strings.TrimSuffix(dst, filepath.Ext(dst))
	err := cloneFile(src, trimmedNewPath)
	if err != nil {
		return fmt.Errorf("cloning file: %w", err)
	}
	f, err := os.OpenFile(trimmedNewPath, os.O_RDWR, 0444) // only enable reading mode as we do not need to write anything
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

	return tmplEx.Execute(f, vj)
	// return nil
}
