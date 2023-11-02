package generate

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"

	"github.com/vision-cli/vision/plugins/vision-plugin-plugin-v1/cmd/initialise"
)

//go:embed all:template
var templateFiles embed.FS

var GenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate the code from templates",
	Long:  "generate code in the template files for the Plugin plugin using the values in the vision.json file",
	RunE:  generateAndCheck,
}

type success struct {
	Success bool `json:"success"`
}

// wraps the run function to determine a success or failed response
func generateAndCheck(cmd *cobra.Command, args []string) error {
	err := run(cmd, args)
	jEnc := json.NewEncoder(os.Stdout)
	if err != nil {
		_ = jEnc.Encode(success{Success: false})
		return fmt.Errorf("generating template: %w", err)
	}

	_, err = exec.Command("go", "mod", "tidy").Output()
	if err != nil {
		return fmt.Errorf("running 'go mod tidy': %w", err)
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

func openVisionJson(vPath string) (*initialise.PluginData, error) {
	f, err := os.OpenFile(vPath, os.O_RDWR, 0444)
	if err != nil {
		return nil, fmt.Errorf("opening config file: %w", err)
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("reading bytes: %w", err)
	}

	type convertConfig struct {
		PluginData initialise.PluginConfig `json:"plugin"`
	}

	var convConf convertConfig
	var jsonData initialise.PluginData

	if err = json.Unmarshal(b, &convConf); err != nil {
		return nil, fmt.Errorf("unmarshalling json: %w", err)
	}

	// convert struct to use correct JSON tag
	jsonData.PluginConfig = convConf.PluginData

	return &jsonData, nil
}

// if path is a directory, just copy it
func cloneDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

// if file isn't template file, just copy it
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

func cloneExecTmpl(src, dst string, vj *initialise.PluginData) error {
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
}
