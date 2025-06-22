package views

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

var _ tea.Model = Blank{}

type Blank struct {
	cont table.Model
	h    int
}

// Init implements tea.Model.
func (b Blank) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (b Blank) Update(tea.Msg) (tea.Model, tea.Cmd) {
	return b, nil
}

// View implements tea.Model.
func (b Blank) View() string {
	return "Loading..."
}
