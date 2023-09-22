package cmd

import (
	_ "embed"
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	initialise "github.com/vision-cli/vision/cmd/init"
	"github.com/vision-cli/vision/execute"
)

// Finds available plugins and initialises them into commands
func init() {
	rootCmd.AddCommand(initialise.RootCmd)
	rootCmd.Flags().AddFlagSet(initVisionFlags())
	pluginMap := findVisionPlugins()
	for plName, plVersion := range pluginMap {
		cobraCmd := &cobra.Command{
			Use:     plName,
			Version: plVersion,
			RunE:    pluginCommand,
		}
		rootCmd.AddCommand(cobraCmd)
	}
}

type jsonInput struct {
	Name string
}

func (ji *jsonInput) String() string {
	b, err := json.Marshal(*ji)
	if err != nil {
		return "error marshalling json"
	}
	return string(b)
}

func (ji *jsonInput) Set(s string) error {
	return json.Unmarshal([]byte(s), ji)
}

func (ji *jsonInput) Type() string {
	return ""
}

var inputFile jsonInput

func initVisionFlags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("vision flag set", 1)
	fs.VarP(&inputFile, "json-input", "j", "parse json into vision")
	return fs
}

var pluginCommand = func(cmd *cobra.Command, args []string) error {
	exe := execute.NewPluginExecutor()
	var arg string
	if len(args) < 1 { // prevents index out of range
		arg = ""
		log.Warnf("No argument provided. Try: \n\t\n vision %v -v", cmd.Use)
	} else {
		arg = args[0]
	}
	exe.RunCommand(cmd.Use, arg)
	return nil
}

//go:embed vision-help.txt
var visionHelp string

var rootCmd = &cobra.Command{
	Use:     "vision",
	Short:   "A developer productivity tool",
	Long:    `Vision is a tool to create microservice platforms and microservice scaffolding code`,
	Example: visionHelp,
	RunE:    inputSelector,
}

var inputSelector = func(cmd *cobra.Command, args []string) error {
	data, err := ioutil.ReadFile(args[1])
	if err != nil {
		return err
	}
	inputFile.Set(string(data))
	log.Info(inputFile.String())
	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
	}
}

// Seaches all paths on the system to find binaries with specific vision formatting and assigns them to a map.
// The formatting is `vision-plugin-[plugin-name]-[version-number]`
func findVisionPlugins() map[string]string {
	var pluginMap = make(map[string]string)

	for _, path := range strings.Split(os.Getenv("PATH"), ":") {
		filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
			if info == nil {
				return nil
			}
			if info.IsDir() {
				return nil
			}
			if strings.HasPrefix(info.Name(), "vision-plugin") {
				pluginSplit := strings.Split(info.Name(), "-")

				pluginMap[pluginSplit[2]] = pluginSplit[3]
			}
			return nil
		})
	}
	return pluginMap
}
