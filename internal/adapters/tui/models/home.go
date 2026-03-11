package models

import (
	"fmt"
	"io"
	"strings"

	"github.com/aboubakary833/codexa/internal/adapters/tui/common"
	"github.com/aboubakary833/codexa/internal/adapters/tui/styles"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var logo = `
 ██████╗ ██████╗ ██████╗ ███████╗██╗  ██╗ █████╗ 
██╔════╝██╔═══██╗██╔══██╗██╔════╝╚██╗██╔╝██╔══██╗
██║     ██║   ██║██║  ██║█████╗   ╚███╔╝ ███████║
██║     ██║   ██║██║  ██║██╔══╝   ██╔██╗ ██╔══██║
╚██████╗╚██████╔╝██████╔╝███████╗██╔╝ ██╗██║  ██║
 ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝
`

var (
	logoStyles            = lipgloss.NewStyle().Foreground(styles.PrimaryColor)
	homeNormalItemStyle   = styles.NormalTitle.Transform(strings.ToUpper).Bold(true)
	homeSelectedItemStyle = styles.SelectedTitle.Transform(strings.ToUpper).Bold(true)
)

// model is the model that render the home(default) view,
// which list the availble view options.
type HomeModel struct {
	width  int
	height int
	list   list.Model
}

func NewHomeModel() HomeModel {
	delegate := common.Delegate{
		RenderFunc: delegateRender,
	}

	items := []list.Item{
		common.ListItem{
			Name: "Browse snippets by category",
			Msg:  LoadTechsMsg{},
		},
		common.ListItem{
			Name: "Search a snippet",
			Msg:  OpenSearchMsg{},
		},
		common.ListItem{
			Name: "About",
			Msg:  OpenAboutMsg{},
		},
	}

	list := common.NewList(
		items, delegate,
		common.WithDefaultListKeyMap,
	)

	list.SetShowTitle(false)
	list.SetShowFilter(false)
	list.SetFilteringEnabled(false)

	list.InfiniteScrolling = true

	list.AdditionalFullHelpKeys = additionalFullHelpKeys

	return HomeModel{
		width:  100,
		height: viewportMaxWidth,
		list:   list,
	}

}

func (m HomeModel) Init() tea.Cmd {
	return tea.SetWindowTitle("Codexa - Home")
}

func (m HomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "q", tea.KeyEsc.String():
			return m, tea.Quit
		}

		m.list, cmd = m.list.Update(msg)
		return m, cmd

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		m.list.SetSize(msg.Width, msg.Height-lipgloss.Height(logo))

		return m, nil
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m HomeModel) View() string {

	if m.height < 16 {
		m.list.SetHeight(m.height)

		return m.list.View()
	}

	if m.width < 50 {
		m.list.SetHeight(m.height)
		return lipgloss.PlaceVertical(
			m.height, lipgloss.Center,
			m.list.View(),
		)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		logoStyles.Render(logo),
		m.list.View(),
	)
}

// delegateRender is the render function of the HomeModel list delegate
func delegateRender(w io.Writer, m list.Model, _ common.Delegate, index int, item list.Item) {
	viewItem, ok := item.(common.ListItem)

	if !ok {
		return
	}

	content := lipgloss.NewStyle().Width(44).Render(viewItem.Name)

	if item.FilterValue() == m.SelectedItem().FilterValue() {
		content = homeSelectedItemStyle.Render(
			content, common.DefaultSelectedItemArrow,
		)
	} else {
		content = homeNormalItemStyle.Render(content)
	}

	fmt.Fprintf(w, "%s", content)
}

// additionalFullHelpKeys return a slice of additional
// keybindings to be added to the help full view of the home view.
func additionalFullHelpKeys() []key.Binding {
	lak := common.NewListAdditionalKeyMap()
	return []key.Binding{
		lak.FullEnter,
		lak.CloseHelp,
		key.NewBinding(
			key.WithKeys("q", "esc"),
			key.WithHelp("q/esc", "quit"),
		),
	}
}
