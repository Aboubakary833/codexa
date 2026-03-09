package models

import (
	"fmt"
	"strings"

	"github.com/aboubakary833/codexa/internal/adapters/tui/common"
	"github.com/aboubakary833/codexa/internal/adapters/tui/styles"
	"github.com/aboubakary833/codexa/internal/domain"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

const viewportMaxWidth = 65

var (
	entryTitleStyle = lipgloss.NewStyle().
			Background(styles.PrimaryColor).
			Foreground(lipgloss.Color("230")).
			Padding(0, 1)

	entryLoadingTextStyle = lipgloss.NewStyle().
				Foreground(styles.DimmedColor).
				Transform(strings.ToUpper).
				Italic(true).Bold(true).
				SetString("loading...")

	entryScrollInfoStyle = lipgloss.NewStyle().
				Background(styles.SecondaryColor).
				Foreground(lipgloss.Color("230")).
				Padding(0, 1)
)

type entryKeyMap struct {
	CursorUp   key.Binding
	CursorDown key.Binding
	Prev       key.Binding
	Quit       key.Binding
	ForceQuit  key.Binding
}

// newKeyMap return a new keyMap for the entry view keybindings.
func newEntryKeyMap() entryKeyMap {
	return entryKeyMap{
		CursorUp: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		CursorDown: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Prev: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "prev"),
		),

		Quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit"),
		),

		ForceQuit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "quit"),
		),
	}
}

// LoadSnippetMsg trigger a given snippet content rendering
type LoadSnippetMsg struct {
	Snippet domain.Snippet
}

// SnippetLoadedMsg is the response message to LoadSnippetMsg.
type SnippetLoadedMsg struct {
	Tech    string
	Topic   string
	Content string
}

// SnippetViewModel represents the model that handle the display
// and navigation of a single entry content.
type SnippetViewModel struct {
	width  int
	height int

	tech    string
	topic   string
	content string

	ready    bool
	viewport viewport.Model

	keyMap entryKeyMap
	help   help.Model
}

func NewSnippetViewModel(tech, topic, content string) SnippetViewModel {

	markdown, err := glamourRender(content)

	// in case where we got error, we ignore it. We would not like
	// to stop snippet rendering because styling failed, do we?
	if err != nil {
		markdown = content
	}

	return SnippetViewModel{
		width:  100,
		height: viewportMaxWidth,

		tech:    tech,
		topic:   topic,
		content: markdown,

		ready: false,

		keyMap: newEntryKeyMap(),
		help:   help.New(),
	}
}

func (m SnippetViewModel) Init() tea.Cmd {
	return nil
}

func (m SnippetViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch {
		case key.Matches(msg, m.keyMap.Prev):
			return m, common.NavigateToPrev

		case key.Matches(msg, m.keyMap.Quit, m.keyMap.ForceQuit):
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		var viewportWidth int

		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		m.width = msg.Width
		m.height = msg.Height
		viewportWidth = min(msg.Width, viewportMaxWidth)

		if !m.ready {
			m.viewport = viewport.New(viewportWidth, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight

			m.viewport.SetContent(m.content)
			m.ready = true
		} else {
			m.viewport.Width = viewportWidth
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

		if m.width < 50 {
			m.viewport.SetHorizontalStep(1)
		} else {
			m.viewport.SetHorizontalStep(0)
		}

		return m, nil
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m SnippetViewModel) ShortHelp() []key.Binding {
	return []key.Binding{
		m.keyMap.CursorUp,
		m.keyMap.CursorDown,
		m.keyMap.Prev,
		m.keyMap.Quit,
	}
}

func (m SnippetViewModel) headerView() string {
	text := fmt.Sprintf("%s:%s", m.tech, m.topic)
	title := entryTitleStyle.Render(strings.ToLower(text))
	barWidth := m.viewport.Width - lipgloss.Width(title)
	headerBar := lipgloss.NewStyle().Width(barWidth).Render()

	return lipgloss.JoinHorizontal(
		lipgloss.Bottom,
		headerBar, title,
	)
}

func (m SnippetViewModel) footerView() string {

	scrollPercentStr := fmt.Sprintf("scroll:%3.f%%", m.viewport.ScrollPercent()*100)
	helpView := m.help.ShortHelpView(m.ShortHelp())
	scrollInfoView := entryScrollInfoStyle.Render(scrollPercentStr)

	if m.width < 50 {
		scrollInfoView = entryScrollInfoStyle.Width(
			lipgloss.Width(helpView),
		).Render(scrollPercentStr)

		return lipgloss.JoinVertical(
			lipgloss.Left,
			helpView, scrollInfoView,
		)
	}

	var whitespaces string
	usedWidth := lipgloss.Width(helpView) + lipgloss.Width(scrollInfoView)
	whitespacesWidth := m.viewport.Width - usedWidth
	if whitespacesWidth > 0 {
		whitespaces = strings.Repeat(" ", whitespacesWidth)
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, helpView, whitespaces, scrollInfoView)
}

func (m SnippetViewModel) View() string {
	var container string

	if m.ready {
		container = m.viewport.View()
	} else {
		viewportHeight := m.viewport.Height
		container = lipgloss.NewStyle().
			Width(m.viewport.Width).Height(viewportHeight).
			Render(lipgloss.Place(
				m.viewport.Width, viewportHeight, lipgloss.Center,
				lipgloss.Center, entryLoadingTextStyle.String(),
			))
	}

	return lipgloss.JoinVertical(
		lipgloss.Left, m.headerView(),
		container, m.footerView(),
	)
}

// glamourRender instantiate a new term renderer and render
// given content with custom stylesheet
func glamourRender(content string) (string, error) {
	styles := glamour.WithStyles(styles.GetRendererStyles())
	renderer, err := glamour.NewTermRenderer(styles)

	if err != nil {
		return "", err
	}

	return renderer.Render(content)
}
