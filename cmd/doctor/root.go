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
		var reasons []string
		info, err := exe.Info()
		// TODO(luke): add "not a string" catch to empty string checks
		switch {
		case err != nil:
			reasons = append(reasons, fmt.Sprintf("%v", err))
		case info.ShortDescription == "":
			reasons = append(reasons, "short description missing")
		case info.LongDescription == "":
			reasons = append(reasons, "log description missing")
		}

		if len(reasons) > 0 {
			invalidPlugins = append(invalidPlugins, ErrInvalidPlugin{
				PluginName: p.Name,
				Command:    "info",
				Reasons:    reasons,
			})
		}

		// command = "init"
		// ini, err := exe.Init()
		// if err != nil {
		// 	reason = fmt.Sprintf("%v", err)
		// 	healthRecord = append(healthRecord, healthLog{
		// 		pluginName:  p.Name,
		// 		command:     command,
		// 		description: reason,
		// 	})
		// } else if ini.Config == "" {
		// 	reason = "config empty"
		// 	healthRecord = append(healthRecord, healthLog{
		// 		pluginName:  p.Name,
		// 		command:     command,
		// 		description: reason,
		// 	})
		// } else if ini.Config == nil {
		// 	reason = "config missing"
		// 	healthRecord = append(healthRecord, healthLog{
		// 		pluginName:  p.Name,
		// 		command:     command,
		// 		description: reason,
		// 	})
		// }

		// command = "version"
		// vers, err := exe.Version()
		// if err != nil {
		// 	reason = fmt.Sprintf("%v", err)
		// 	healthRecord = append(healthRecord, healthLog{
		// 		pluginName:  p.Name,
		// 		command:     command,
		// 		description: reason,
		// 	})
		// } else if vers.SemVer == "" {
		// 	reason = "semantic version missing"
		// 	healthRecord = append(healthRecord, healthLog{
		// 		pluginName:  p.Name,
		// 		command:     command,
		// 		description: reason,
		// 	})
		// }
	}
	// printTable(healthRecord)
	// healthRecordPrinter(healthRecord)
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
