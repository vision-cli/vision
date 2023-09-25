package styles

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	DoctorInfoStyle = lipgloss.NewStyle().
			MarginLeft(1).
			MarginRight(1).
			Padding(0, 1).
			Bold(true).
			SetString("Health Check for plugin")

	DoctorPluginNameStyle = lipgloss.NewStyle().
				MarginLeft(0).
				MarginRight(5).
				Padding(0, 1).
				Bold(true).
				Foreground(lipgloss.Color("#fcc783"))

	DoctorTableStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("240"))
)

type model struct {
	table table.Model
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, tea.Quit
}

func (m model) View() string {
	return DoctorTableStyle.Render(m.table.View()) + "\n"
}

func ShowTable(columns []table.Column, rows []table.Row) {

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		// table.WithFocused(true),
		table.WithHeight(len(rows)),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	s.Cell.Foreground(lipgloss.Color("231"))

	t.SetStyles(s)

	m := model{t}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

}
