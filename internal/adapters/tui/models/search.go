package models

import (
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/aboubakary833/codexa/internal/adapters/tui/common"
	"github.com/aboubakary833/codexa/internal/adapters/tui/helpers"
	"github.com/aboubakary833/codexa/internal/adapters/tui/styles"
	"github.com/aboubakary833/codexa/internal/domain"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	searchInputStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				Padding(0, 1).MaxWidth(50)

	searchInputFocusStyle = searchInputStyle.
				BorderForeground(styles.PrimaryColor)

	searchInputBlurStyle = searchInputStyle.
				BorderForeground(styles.TextColor)

	emptySearchListStyle = lipgloss.NewStyle().
				SetString("No available snippet entry found").
				MarginLeft(1)

	noMatchingEntryStyle = lipgloss.NewStyle().
				SetString("No matching snippet entry found").
				MarginLeft(1)
)

// OpenSearchMsg request the opening of the search view.
// This will trigger the display of the SearchModel view.
type OpenSearchMsg struct{}

// SearchRequestMsg send a search request with a query.
type SearchRequestMsg struct {
	Query string
}

// SearchResponseMsg is the result to SearchRequestMsg.
// It is handled by the SearchModel that will update the list items set.
type SearchResponseMsg struct {
	Snippets []domain.Snippet
}

// These two focus message are used to tweak the focus
// between input and list delegate.
type FocusInputMsg struct{}
type FocusListMsg struct{}

type emptyListKeyMap struct {
	Prev key.Binding
	Quit key.Binding
}

func (km emptyListKeyMap) ToBindings() []key.Binding {
	return []key.Binding{km.Prev, km.Quit}
}

type SearchModel struct {
	width  int
	height int

	input           textinput.Model
	listDelegate    *common.Delegate
	list            list.Model
	emptyListKeyMap emptyListKeyMap
}

func NewSearchModel() SearchModel {
	searchInput := textinput.New()

	searchInput.Placeholder = "ex: go/context-timeout, php arrays..."
	searchInput.Focus()

	delegate := &common.Delegate{
		RenderFunc: searchListDelegateRender,
		UpdateFunc: searchListDelegateUpdate,
		Focused:    false,
	}

	resultList := common.NewList(
		[]list.Item{}, delegate,
		common.WithDefaultListKeyMap,
	)
	resultList.SetDelegate(delegate)

	resultList.SetShowTitle(false)
	resultList.SetShowFilter(false)
	resultList.SetFilteringEnabled(false)
	resultList.InfiniteScrolling = false

	return SearchModel{
		width: viewportMaxWidth,
		height: 100,
		input:        searchInput,
		list:         resultList,
		listDelegate: delegate,
		emptyListKeyMap: emptyListKeyMap{
			Prev: key.NewBinding(
				key.WithKeys("esc"),
				key.WithHelp("esc", "prev view"),
			),

			Quit: key.NewBinding(
				key.WithKeys("ctrl+c"),
				key.WithHelp("ctrl+c", "quit"),
			),
		},
	}
}

func (m SearchModel) Init() tea.Cmd {
	return doSearch("")
}

func (m SearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp, tea.KeyPgUp:
			if m.listDelegate.Focused {
				if m.list.GlobalIndex() == 0 {
					return m, focusInput()
				}
				m.list, cmd = m.list.Update(msg)
				return m, cmd
			}

		case tea.KeyEsc:
			if m.input.Focused() {
				return m, common.NavigateToPrev
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		if msg.Width >= 50 {
			m.input.Width = 43
		} else {
			m.input.Width = msg.Width - 10
		}
		m.list.SetSize(msg.Width, msg.Height-lipgloss.Height(m.input.View())-2)

		return m, nil

	case FocusInputMsg:
		if !m.input.Focused() {
			m.listDelegate.Focused = false
			cmd = m.input.Focus()

			return m, cmd
		}

		return m, nil

	case FocusListMsg:
		if !m.listDelegate.Focused {
			m.listDelegate.Focused = true
		}
		m.input.Blur()

		return m, nil

	case SearchResponseMsg:
		items := makeItemsFromSearchResult(msg.Snippets)
		cmd = m.list.SetItems(items)
		m.list.Select(0)

		return m, cmd
	}

	if !m.input.Focused() {
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}

	return m.handleInputMsg(msg)
}

func (m SearchModel) listEmpty() bool {
	return len(m.list.Items()) == 0
}

func (m SearchModel) listView() string {
	if m.listEmpty() {
		helpView := m.list.Help.ShortHelpView(
			m.emptyListKeyMap.ToBindings(),
		)
		style := emptySearchListStyle

		if m.input.Value() != "" {
			style = noMatchingEntryStyle
		}

		return lipgloss.JoinVertical(
			lipgloss.Left,
			style.Width(48).
				Height(m.list.Height()-lipgloss.Height(helpView)).
				Render(), helpView,
		)
	}

	return m.list.View()
}

func (m SearchModel) View() string {
	var (
		input      string
		inputWidth = m.input.Width + 7
	)

	if m.width >= 55 {
		inputWidth = inputWidth - 2
	}

	if m.input.Focused() {
		input = searchInputFocusStyle.Width(inputWidth).
			Render(m.input.View())
	} else {
		input = searchInputBlurStyle.Width(inputWidth).
			Render(m.input.View())
	}

	return lipgloss.JoinVertical(
		lipgloss.Left, input,
		m.listView(),
	)
}

// handleInputMsg handle the message for the input model
func (m SearchModel) handleInputMsg(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		onBlurKeys := []tea.KeyType{tea.KeyTab, tea.KeyDown, tea.KeyPgDown}

		if slices.Contains(onBlurKeys, msg.Type) && !m.listEmpty() {
			return m, focusList()
		}

		if key.Matches(msg, m.emptyListKeyMap.Quit) {
			return m, tea.Quit
		}

		if msg.Type == tea.KeyEnter && !m.listEmpty() {
			cmd = doSearch(m.input.Value())
			cmds = append(cmds, cmd, focusList())

			return m, tea.Batch(cmds...)
		}

		m.input, cmd = m.input.Update(msg)
		searchCmd := doSearch(m.input.Value())

		cmds = append(cmds, cmd, searchCmd)

		return m, tea.Batch(cmds...)
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

// doSearch make a new search request
func doSearch(query string) tea.Cmd {
	return func() tea.Msg {
		return SearchRequestMsg{
			Query: query,
		}
	}
}

func focusInput() tea.Cmd {
	return func() tea.Msg {
		return FocusInputMsg{}
	}
}

func focusList() tea.Cmd {
	return func() tea.Msg {
		return FocusListMsg{}
	}
}

func searchListDelegateUpdate(msg tea.Msg, m *list.Model, d *common.Delegate) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:

		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return tea.Quit
		}

		switch msg.Type {

		case tea.KeyEnter:
			item, ok := m.SelectedItem().(searchItem)

			if !ok {
				return nil
			}

			return item.LoadEntry()

		case tea.KeyEsc, tea.KeyShiftTab:
			if d.Focused {
				return focusInput()
			}
		}

	}

	return nil
}

// searchListDelegateRender is the search result list delegate render function
func searchListDelegateRender(w io.Writer, m list.Model, d common.Delegate, index int, item list.Item) {
	searchItem, ok := item.(searchItem)

	if !ok {
		return
	}
	var contentWidth int

	if m.Width() >= 50 {
		contentWidth = 45
	} else {
		contentWidth = m.Width() - 5
	}

	content := lipgloss.NewStyle().Width(contentWidth).Render(searchItem.Title())

	if searchItem.FilterValue() == m.SelectedItem().FilterValue() && d.Focused {
		content = styles.SelectedTitle.MaxWidth(50).
			BorderStyle(lipgloss.NormalBorder()).Render(
			content, common.DefaultSelectedItemArrow,
		)
	} else {
		content = styles.NormalTitle.Render(content)
	}

	fmt.Fprintf(w, "%s", content)
}

type searchItem domain.Snippet

func (i searchItem) Title() string {
	return fmt.Sprintf("%s %s", i.TechID, strings.ToLower(i.Topic))
}

func (i searchItem) FilterValue() string {
	return i.ID
}

func (i searchItem) LoadEntry() tea.Cmd {
	return func() tea.Msg {
		return LoadSnippetMsg{
			Snippet: domain.Snippet{
				ID:       i.ID,
				TechID:   i.TechID,
				Topic:    i.Topic,
				Filepath: i.Filepath,
			},
		}
	}
}

func makeItemsFromSearchResult(snippets []domain.Snippet) []list.Item {
	return helpers.MakeItems(snippets, func(snippet domain.Snippet) list.Item {
		return searchItem(snippet)
	})
}
