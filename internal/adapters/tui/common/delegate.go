package common

import (
	"fmt"
	"io"

	"github.com/aboubakary833/codexa/internal/adapters/tui/styles"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const DefaultSelectedItemArrow = "»"

type DelegateRenderFunc func(w io.Writer, m list.Model, d Delegate, index int, item list.Item)
type DelegateUpdateFunc func(msg tea.Msg, m *list.Model, d *Delegate) tea.Cmd

type Delegate struct {
	// These methods bypass the default delegate methods
	UpdateFunc DelegateUpdateFunc
	RenderFunc DelegateRenderFunc

	// The arrow character of the selected item. Default to "»"
	SelectedItemArrow string
	Focused           bool
}

// NewDelegate return a new Delegate with UpdateFunc
// and RenderFunc not set.
func NewDelegate() Delegate {
	return Delegate{
		SelectedItemArrow: DefaultSelectedItemArrow,
		Focused:           true,
	}
}

func (d Delegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {

	if d.UpdateFunc == nil {
		return d.defaultUpdateFunc(msg, m)
	}

	return d.UpdateFunc(msg, m, &d)
}

func (d Delegate) Render(w io.Writer, m list.Model, index int, item list.Item) {

	if d.RenderFunc == nil {
		d.defaultRenderFunc(w, m, index, item)
		return
	}

	d.RenderFunc(w, m, d, index, item)
}

func (d Delegate) Height() int {
	return 1
}

func (d Delegate) Spacing() int {
	return 1
}

// DelegateItem specifies the contract expected by this Delegate.
// Lists that rely on this Delegate must implement this interface,
// unless they provide their own custom RenderFunc.
type DelegateItem interface {
	list.Item
	Load() tea.Cmd
	Title() string
}

// defaultRenderFunc is the fallback render function that will be run if no custom
// RenderFunc has been provided for the current Delegate struct.
func (d Delegate) defaultRenderFunc(w io.Writer, m list.Model, _ int, item list.Item) {
	delegateItem, ok := item.(DelegateItem)

	if !ok {
		return
	}
	var contentWidth int

	if m.Width() > 52 {
		contentWidth = 48
	} else {
		contentWidth = m.Width() - 5
	}

	content := lipgloss.NewStyle().Width(contentWidth).Render(delegateItem.Title())

	if delegateItem.FilterValue() == m.SelectedItem().FilterValue() && d.Focused {
		content = styles.SelectedTitle.Render(content, d.SelectedItemArrow)
	} else {
		content = styles.NormalTitle.Render(content)
	}

	fmt.Fprintf(w, "%s", content)
}

// defaultUpdateFunc is the delegate fallback update function
func (d Delegate) defaultUpdateFunc(msg tea.Msg, m *list.Model) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyEnter.String():
			defaultItem, ok := m.SelectedItem().(DelegateItem)

			if !ok {
				return nil
			}
			return defaultItem.Load()

		case "q", "ctrl+c":
			return tea.Quit
		}

	}

	return nil
}
