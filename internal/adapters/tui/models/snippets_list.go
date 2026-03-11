package models

import (
	"fmt"

	"github.com/aboubakary833/codexa/internal/adapters/tui/common"
	"github.com/aboubakary833/codexa/internal/adapters/tui/helpers"
	"github.com/aboubakary833/codexa/internal/domain"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// LoadTechSnippetsMsg request the controller to fetch and return
// a specific category entries
type LoadTechSnippetsMsg struct {
	Tech domain.Tech
}

// SnippetsLoadedMsg is the successful result of LoadTechSnippetsMsg.
// It also load to the display of SnippetsModel view with the list of snippets.
type TechSnippetsLoadedMsg struct {
	Tech     domain.Tech
	Snippets []domain.Snippet
}

// SnippetsListModel render a list of a specific tech category snippets
type SnippetsListModel struct {
	width            int
	height           int

	tech			 domain.Tech
	list             list.Model
	additionalKeyMap common.ListAdditionalKeyMap
}

func NewEntries(tech domain.Tech, snippets []domain.Snippet) SnippetsListModel {

	items := helpers.MakeItems(snippets, func(i domain.Snippet) list.Item {
		return common.ListItem{
			Name: i.Topic,
			Msg: LoadSnippetMsg{
				Snippet: i,
			},
		}
	})
	delegate := common.NewDelegate()
	listkeyMap := common.NewListAdditionalKeyMap()

	list := common.NewList(items, delegate)
	list.Title = fmt.Sprintf("%s snippets entries", tech.Name)

	list.AdditionalFullHelpKeys = listkeyMap.FullKeys
	list.AdditionalShortHelpKeys = listkeyMap.ShortKeys

	return SnippetsListModel{
		width: viewportMaxWidth,
		height: 100,
		tech: tech,
		list:             list,
		additionalKeyMap: listkeyMap,
	}
}

func (m SnippetsListModel) Init() tea.Cmd {
	title := fmt.Sprintf("%s snippets", m.tech.Name)
	return tea.SetWindowTitle(title)
}

func (m SnippetsListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
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

	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m SnippetsListModel) View() string {
	if !helpers.IsListUnfiltered(m.list) {
		m.list.AdditionalFullHelpKeys = m.additionalKeyMap.FullKeysWithoutPrev
	}

	return m.list.View()
}
