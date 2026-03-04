package models

import (
	"github.com/aboubakary833/codexa/internal/adapters/tui/common"
	"github.com/aboubakary833/codexa/internal/adapters/tui/helpers"
	"github.com/aboubakary833/codexa/internal/adapters/tui/styles"
	"github.com/aboubakary833/codexa/internal/domain"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	cmdStyle = lipgloss.NewStyle().
			Foreground(styles.SecondaryColor).
			SetString("codexa pull")
)

// LoadTechsMsg request the listing of tech categories
type LoadTechsMsg struct{}

// TechsLoadedMsg return a slice of tech categories.
// It is the message that lead to the display of tech categories list.
type TechsLoadedMsg struct {
	Techs []domain.Tech
}

// TechsListModel list the available selectable categories.
type TechsListModel struct {
	width int
	height int
	list list.Model
	additionalKeyMap common.ListAdditionalKeyMap
}

func NewTechsListModel(techs []domain.Tech) TechsListModel {

	items := helpers.MakeItems(techs, func(tech domain.Tech) list.Item {
		return common.ListItem{
			Name: tech.Name,
			Msg:  LoadTechSnippetsMsg{Tech: tech},
		}
	})

	delegate := common.NewDelegate()
	listkeyMap := common.NewListAdditionalKeyMap()

	list := common.NewList(items, delegate)
	list.Title = "Available tech categories"

	list.AdditionalFullHelpKeys = listkeyMap.FullKeys
	list.AdditionalShortHelpKeys = listkeyMap.ShortKeys

	return TechsListModel{
		list: list,
		additionalKeyMap: listkeyMap,
	}
}

func (m TechsListModel) Init() tea.Cmd {
	return nil
}

func (m TechsListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	
	case tea.KeyMsg:
		if msg.Type == tea.KeyEsc && helpers.IsListUnfiltered(m.list) {
			return m, common.NavigateToPrev
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		listWidth := min(msg.Width, 58)
		m.list.SetSize(listWidth, msg.Height)
		m.list.Styles.Title = m.list.Styles.Title.Width(msg.Width)

		return m, nil
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m TechsListModel) listEmpty() bool {
	return len(m.list.Items()) == 0
}

func (m TechsListModel) View() string {
	if m.listEmpty() {
		return m.emptyListView()
	}

	if !helpers.IsListUnfiltered(m.list) {
		m.list.AdditionalFullHelpKeys = m.additionalKeyMap.FullKeysWithoutPrev
	}

	return m.list.View()
}

func (m TechsListModel) emptyListView() string {
	titleView := styles.ListTitleStyle.
		Width(m.list.Width()).Render(
			m.list.Title,
		)
	helpView := m.list.Help.ShortHelpView(
		[]key.Binding{
			m.additionalKeyMap.Prev,
			m.additionalKeyMap.Quit,
		},
	)
	verticalMargin := lipgloss.Height(titleView) + lipgloss.Height(helpView)

	emptyListView := lipgloss.NewStyle().Width(m.list.Width()).
		Height(m.list.Height() - verticalMargin - 1).
		MarginTop(1).Render(
			"No snippet tech category found.",
			"\nRun", cmdStyle.Render(), "to fetch a tech category of snippets.",
		)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		titleView,
		emptyListView,
		helpView,
	)
}
