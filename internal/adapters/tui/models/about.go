package models

import (
	"fmt"
	"strings"

	"github.com/aboubakary833/codexa/internal/adapters/tui/common"
	"github.com/aboubakary833/codexa/internal/adapters/tui/styles"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	title             = "About Codexa"
	content           = "is a terminal-based application designed to help developers quickly access concise, practical code snippets and common development patterns,without the verbosity of traditional documentation."
	contribution      = "To contribute, visit one of the following urls:"
	registryUrl       = "github.com/aboubakary833/cx-registry"
	sourceCodeRepoUrl = "github.com/aboubakary833/codexa"
)

const (
	defaultWidth = 60
	defaultPadding = 1
)

var (
	aboutTitleStyle = lipgloss.NewStyle().Foreground(styles.PrimaryColor).
			Border(lipgloss.BlockBorder(), false, false, false, true).
			BorderForeground(styles.PrimaryColor).Padding(0, 1).MarginBottom(1).
			Transform(strings.ToUpper).Bold(true)

	urlStyle = lipgloss.NewStyle().Foreground(styles.PrimaryColor).
			Underline(true).Bold(true)
)

// OpenAboutMsg is used to request
// the display of the AboutModel view
type OpenAboutMsg struct{}

// model renders the about view
type AboutModel struct {
	width  int
	height int
	help   help.Model
}

func NewAboutModel() AboutModel {
	return AboutModel{
		width:  100,
		height: viewportMaxWidth,
		help:   help.New(),
	}
}

func (m AboutModel) Init() tea.Cmd {
	return nil
}

func (m AboutModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.KeyMsg:
		if msg.Type == tea.KeyEsc {
			return m, common.NavigateToPrev
		}

		if msg.String() == "q" || msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		return m, nil
	}

	return m, nil
}

func (m AboutModel) helpView() string {
	keys := []key.Binding{
		key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "go back to home"),
		),
		key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}

	return lipgloss.NewStyle().PaddingLeft(defaultPadding).Render(
		m.help.ShortHelpView(keys),
	)
}

func (m AboutModel) View() string {
	var sections []string

	sections = append(sections, aboutTitleStyle.Render(title))
	availableHeight := m.height - lipgloss.Height(m.helpView())

	containerStyle := lipgloss.NewStyle().Width(defaultWidth).Height(availableHeight).Padding(defaultPadding)
	styledName := lipgloss.NewStyle().Foreground(styles.PrimaryColor).Render("Codexa")
	styledDesc := lipgloss.NewStyle().Foreground(styles.TextColor).Render(content)

	desc := fmt.Sprintf("%s %s", styledName, styledDesc)
	contribute := fmt.Sprintf(
		"\n%s\n\n%s\n%s", contribution,
		urlStyle.Render(sourceCodeRepoUrl),
		urlStyle.Render(registryUrl),
	)

	sections = append(sections, desc, contribute)

	container := containerStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left, sections...),
	)

	return fmt.Sprintf("%s\n%s", container, m.helpView())
}
