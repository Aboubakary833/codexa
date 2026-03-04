package common

import (
	"github.com/aboubakary833/codexa/internal/adapters/tui/styles"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type ListOption func (list.Model) (list.Model)

func NewList(items []list.Item, delegate list.ItemDelegate, options... ListOption) list.Model {
	list := list.New(items, delegate, 0, 0)

	list.Styles.Title = styles.ListTitleStyle
	list.Styles.TitleBar = styles.ListTitleBarStyle
	list.FilterInput.Cursor.Style = styles.ListFilterCursorStyle
	list.FilterInput.PromptStyle = styles.ListFilterPromptStyle
	list.InfiniteScrolling = true

	list.SetShowStatusBar(false)

	list.DisableQuitKeybindings()
	list.KeyMap.CloseFullHelp.Unbind()

	for _, option := range options {
		list = option(list)
	}

	return list
}

// ListItem implements the delegate DefaultItem
type ListItem struct {
	Name string
	Sub  string
	Msg  tea.Msg
}

func (i ListItem) Title() string {
	return i.Name
}

func (i ListItem) FilterValue() string {
	return i.Name
}

func (i ListItem) Load() tea.Cmd {
	return func() tea.Msg {
		return i.Msg
	}
}

// WithDefaultListKeyMap is an option that bind default additional keymap to the list model
func WithDefaultListKeyMap(m list.Model) list.Model {
	lak := NewListAdditionalKeyMap()

	m.AdditionalFullHelpKeys = lak.FullKeys
	m.AdditionalShortHelpKeys = lak.ShortKeys

	return m
}

type ListAdditionalKeyMap struct {
	ShortEnter key.Binding
	FullEnter  key.Binding
	Prev       key.Binding
	CloseHelp  key.Binding
	Quit       key.Binding
}

func NewListAdditionalKeyMap() ListAdditionalKeyMap {
	return ListAdditionalKeyMap{
		ShortEnter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "goto"),
		),

		FullEnter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "go to view"),
		),

		Prev: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "prev view"),
		),
		
		CloseHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "close help"),
		),
		
		Quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit"),
		),

	}
}

// ShortKeys return a slice of additional
// keybindings to be added to the help short view keys slice
func (lak ListAdditionalKeyMap) ShortKeys() []key.Binding {
	return []key.Binding{lak.ShortEnter}
}

// AdditionalFullHelpKeys return a slice of additional
// keybindings to be added to the help full view keys slice
func (lak ListAdditionalKeyMap) FullKeys() []key.Binding {
	return []key.Binding{
		lak.Prev,
		lak.FullEnter,
		lak.CloseHelp,
		lak.Quit,
	}
}

func (lak ListAdditionalKeyMap) FullKeysWithoutPrev() []key.Binding {
	return []key.Binding{
		lak.FullEnter,
		lak.CloseHelp,
		lak.Quit,
	}
}
