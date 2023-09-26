package doctor

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"

	"github.com/vision-cli/vision/internal/plugin"
	"github.com/vision-cli/vision/styles"
)

var RootCmd = &cobra.Command{
	Use:   "doctor",
	Short: "check status of plugins",
	Long:  "Check the status of all plugins and their commands. If a command fails, it will log an error.",
	RunE:  doctorCommand,
}

var healthRecord []healthLog

type healthLog struct {
	pluginName  string
	command     string
	description string
}

// Doctor command looks through available plugins and checks for plugin health.
// If plugins commands are missing or incomplete, doctor returns them as faulty with a reason and prints them out.
var doctorCommand = func(cmd *cobra.Command, args []string) error {
	plugins := plugin.Find()

	for _, plug := range plugins {
		log.Infof("The plugins available are: %v", plug.Name)
	}

	for _, plug := range plugins {
		// call each of the built in commands
		exe := plugin.NewExecutor(plug.FullPath)
		var reason string
		info, err := exe.Info()
		// TODO(luke): add "not a string" catch to empty string checks
		command := "info"
		if err != nil {
			reason = fmt.Sprintf("%v", err)
			healthRecord = append(healthRecord, healthLog{
				pluginName:  plug.Name,
				command:     command,
				description: reason,
			})
		} else if info.ShortDescription == "" {
			reason = "short description missing"
			healthRecord = append(healthRecord, healthLog{
				pluginName:  plug.Name,
				command:     command,
				description: reason,
			})
		} else if info.LongDescription == "" {
			reason = "long description missing"
			healthRecord = append(healthRecord, healthLog{
				pluginName:  plug.Name,
				command:     command,
				description: reason,
			})
		}

		command = "init"
		ini, err := exe.Init()
		if err != nil {
			reason = fmt.Sprintf("%v", err)
			healthRecord = append(healthRecord, healthLog{
				pluginName:  plug.Name,
				command:     command,
				description: reason,
			})
		} else if ini.Config == "" {
			reason = "config empty"
			healthRecord = append(healthRecord, healthLog{
				pluginName:  plug.Name,
				command:     command,
				description: reason,
			})
		} else if ini.Config == nil {
			reason = "config missing"
			healthRecord = append(healthRecord, healthLog{
				pluginName:  plug.Name,
				command:     command,
				description: reason,
			})
		}

		command = "version"
		vers, err := exe.Version()
		if err != nil {
			reason = fmt.Sprintf("%v", err)
			healthRecord = append(healthRecord, healthLog{
				pluginName:  plug.Name,
				command:     command,
				description: reason,
			})
		} else if vers.SemVer == "" {
			reason = "semantic version missing"
			healthRecord = append(healthRecord, healthLog{
				pluginName:  plug.Name,
				command:     command,
				description: reason,
			})
		}
	}
	printTable()
	// healthRecordPrinter(healthRecord)
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

func printTable() {
	columns := []table.Column{
		{Title: "Plugin", Width: 10},
		{Title: "Command", Width: 10},
		{Title: "Fault", Width: 27},
	}
	rows := []table.Row{}
	// TODO(genevieve): fix repeated plugin names
	for n, log := range healthRecord {
		rows = append(rows, table.Row{log.pluginName, log.command, log.description})
		if n+1 <= len(healthRecord)-1 && (healthRecord[n].pluginName != healthRecord[n+1].pluginName) {
			rows = append(rows, table.Row{"\n"})
		}
	}

	styles.ShowTable(columns, rows)
}
