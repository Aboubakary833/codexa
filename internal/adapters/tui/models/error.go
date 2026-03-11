package models

import (
	"fmt"

	"github.com/aboubakary833/codexa/internal/adapters/tui/common"
	"github.com/aboubakary833/codexa/internal/adapters/tui/styles"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var errLogo = `

███████╗██████╗ ██████╗  ██████╗ ██████╗ 
██╔════╝██╔══██╗██╔══██╗██╔═══██╗██╔══██╗
█████╗  ██████╔╝██████╔╝██║   ██║██████╔╝
██╔══╝  ██╔══██╗██╔══██╗██║   ██║██╔══██╗
███████╗██║  ██║██║  ██║╚██████╔╝██║  ██║
╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═╝
`

var errMessageStyle = lipgloss.NewStyle().Foreground(styles.ErrorColor).MarginTop(1)

type ErrorModel struct {
	width  int
	height int

	message string
	help    help.Model
}

func NewErrorModel(message string) ErrorModel {
	return ErrorModel{
		message: message,
		width:   viewportMaxWidth,
		height:  100,
		help: help.New(),
	}
}

func (m ErrorModel) Init() tea.Cmd {
	return tea.SetWindowTitle("Codexa - Error")
}

func (m ErrorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = min(msg.Width, viewportMaxWidth)
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		if msg.Type == tea.KeyEsc {
			return m, common.NavigateToPrev
		}

		if msg.String() == "q" || msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m ErrorModel) helpView() string {
	keys := []key.Binding{
		key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "prev view"),
		),

		key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit"),
		),
	}

	return m.help.ShortHelpView(keys)
}

func (m ErrorModel) View() string {
	helpView := m.helpView()

	if m.height < 16 || m.width < 50 {
		return lipgloss.JoinVertical(
			lipgloss.Left, errMessageStyle.Width(m.width).Height(
				m.height-lipgloss.Height(helpView),
			).Render(m.message), helpView,
		)
	}

	logoView := logoStyles.Width(m.width).Render(errLogo)
	msgVerticalMargin := lipgloss.Height(logoView) + lipgloss.Height(helpView)
	msgView := errMessageStyle.Width(m.width).Height(m.height - msgVerticalMargin).Render(m.message)

	return fmt.Sprintf("%s\n%s\n%s", logoView, msgView, helpView)
}
