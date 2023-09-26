package doctor

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	"github.com/spf13/cobra"

	"github.com/vision-cli/vision/internal/plugin"
	"github.com/vision-cli/vision/styles"
)

var RootCmd = &cobra.Command{
	Use:   "doctor",
	Short: "check status of plugins",
	Long:  "Check the status of all plugins and their commands. If a command fails, it will log an error.",
	RunE:  doctorCmd,
}

var healthRecord []healthLog

type healthLog struct {
	pluginName  string
	command     string
	description string
}

type ErrInvalidPlugin struct {
	PluginName string
	Command    string
	Reasons    []string
}

func (e ErrInvalidPlugin) Error() string {
	out := ""
	out += fmt.Sprintf("Plugin: %s\n", e.PluginName)
	out += fmt.Sprintf("Command: %s\n\n", e.Command)
	for _, reason := range e.Reasons {
		out += fmt.Sprintf("    ->  %s\n", reason)
	}
	return out
}

// doctorCmd looks through available plugins and checks for plugin health.
// If plugins commands are missing or incomplete, doctor returns them as faulty with a reason and prints them out.
var doctorCmd = func(cmd *cobra.Command, args []string) error {
	plugins := plugin.Find()

	var invalidPlugins []error

	for _, p := range plugins {
		// call each of the built in commands
		exe := plugin.NewExecutor(p.FullPath)
		var infoReasons []string
		info, err := exe.Info()
		switch {
		case err != nil:
			infoReasons = append(infoReasons, fmt.Sprintf("%v", err))
		case info.ShortDescription == "":
			infoReasons = append(infoReasons, "short description missing")
		case info.LongDescription == "":
			infoReasons = append(infoReasons, "long description missing")
		}

		if len(infoReasons) > 0 {
			invalidPlugins = append(invalidPlugins, ErrInvalidPlugin{
				PluginName: p.Name,
				Command:    "info",
				Reasons:    infoReasons,
			})
		}

		ini, err := exe.Init()
		var initReasons []string
		switch {
		case err != nil:
			initReasons = append(initReasons, fmt.Sprintf("%v", err))
		case ini.Config == "":
			initReasons = append(initReasons, "config empty")
		case ini.Config == nil:
			initReasons = append(initReasons, "config missing")
		}

		if len(initReasons) > 0 {
			invalidPlugins = append(invalidPlugins, ErrInvalidPlugin{
				PluginName: p.Name,
				Command:    "init",
				Reasons:    initReasons,
			})
		}

		vers, err := exe.Version()
		var versReasons []string
		switch {
		case err != nil:
			versReasons = append(versReasons, fmt.Sprintf("%v", err))
		case vers.SemVer == "":
			versReasons = append(versReasons, "semantic version empty")
		}
	}

	for _, ip := range invalidPlugins {
		fmt.Println(ip.Error())
	}
	return nil
}

// func addToHealthRecord(pluginName string, command string, reason string) {
// 	healthCheck := fmt.Sprintf("%v is faulty. Reason: %v", command, reason)
// 	healthRecord = append(healthRecord, pluginName, healthCheck)
// }

// func healthRecordPrinter(healthRecord []healthLog) {
// 	var curPlugin string
// 	for n, hr := range healthRecord {
// 		// where indexes are even and different to the previous plugin, print out the plugin name
// 		if n == 0 || n%2 == 0 {
// 			if curPlugin != healthRecord[n] {
// 				curPlugin = healthRecord[n]
// 				fmt.Println(styles.DoctorInfoStyle.String() + styles.DoctorPluginNameStyle.Render(strings.ToUpper(curPlugin)))
// 				// fmt.Printf("\nDetails for plugin: %v\n", curPlugin)
// 			}
// 		} else {
// 			log.Warn(hr)
// 		}
// 	}
// 	printTable()
// }

func printTable(logs []healthLog) {
	columns := []table.Column{
		{Title: "Plugin", Width: 10},
		{Title: "Command", Width: 10},
		{Title: "Fault", Width: 27},
	}
	rows := []table.Row{}
	// TODO(genevieve): fix repeated plugin names
	for n, log := range logs {
		rows = append(rows, table.Row{log.pluginName, log.command, log.description})
		if n+1 <= len(healthRecord)-1 && (healthRecord[n].pluginName != healthRecord[n+1].pluginName) {
			rows = append(rows, table.Row{"\n"})
		}
	}

	styles.ShowTable(columns, rows)
}
