package common

import tea "github.com/charmbracelet/bubbletea"

// NavigateToPrevMsg is responsible for triggering
// navigation to previous view
type NavigateToPrevMsg struct {}

// NavigateToPrev is a cmd for going to prev view
var NavigateToPrev = func () tea.Msg {
	return NavigateToPrevMsg{}
}
