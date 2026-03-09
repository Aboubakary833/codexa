package tui

import (
	"log/slog"

	"github.com/aboubakary833/codexa/internal/adapters/tui/common"
	"github.com/aboubakary833/codexa/internal/adapters/tui/models"
	"github.com/aboubakary833/codexa/internal/ports"
	tea "github.com/charmbracelet/bubbletea"
)

// ErrorMsg report when an unexpected error occur
type ErrorMsg struct{
	Error error
}

// rootModel is the TUI entry point. It is the model responsible
// of the inter-model navigation and error-handling when browsing.
type rootModel struct {
	controller   controller
	currentModel tea.Model
	history      *history
}

func New(app ports.Application, initialModel tea.Model, logger *slog.Logger) rootModel {
	controller := newController(app, logger)

	if initialModel == nil {
		initialModel = models.NewHomeModel()
	}

	return rootModel{
		controller:   controller,
		currentModel: initialModel,
		history: &history{
			stack: []tea.Model{},
		},
	}
}

func (m rootModel) Init() tea.Cmd {
	if m.currentModel != nil {
		return m.currentModel.Init()
	}
	return nil
}

func (m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	// Handling current model update messages
	case models.OpenSearchMsg:
		return m.SetCurrentModel(models.NewSearchModel())

	case models.OpenAboutMsg:
		return m.SetCurrentModel(models.NewAboutModel())

	case models.TechsLoadedMsg:
		return m.SetCurrentModel(models.NewTechsListModel(msg.Techs))

	case models.TechSnippetsLoadedMsg:
		return m.SetCurrentModel(models.NewEntries(msg.Tech, msg.Snippets))

	case models.SnippetLoadedMsg:
		return m.SetCurrentModel(models.NewSnippetViewModel(msg.Tech, msg.Topic, msg.Content))

	// Handling controller-request related messages
	case models.LoadTechsMsg:
		return m, m.controller.loadTechCategories()

	case models.LoadTechSnippetsMsg:
		return m, m.controller.loadTechSnippets(msg.Tech)

	case models.LoadSnippetMsg:
		return m, m.controller.loadSnippetContent(msg.Snippet)

	case models.SearchRequestMsg:
		return m, m.controller.search(msg.Query)

	case common.NavigateToPrevMsg:
		return m.SetPreviousModel()

	case ErrorMsg:
		errMsgString := msg.Error.Error()
		return m.SetCurrentModel(models.NewErrorModel(errMsgString))
	}

	var cmd tea.Cmd
	m.currentModel, cmd = m.currentModel.Update(msg)

	return m, cmd
}

// SetCurrentModel update the rootModel to set the new active child model that
// should be rendered. This method return the updated rootModel and WindowSizeMsg cmd
func (m rootModel) SetCurrentModel(newModel tea.Model) (tea.Model, tea.Cmd) {
	if m.currentModel != nil {
		m.history.Push(m.currentModel)
	}
	m.currentModel = newModel

	initCmd := m.currentModel.Init()
	windowSizeCmd := tea.WindowSize()

	return m, tea.Sequence(initCmd, windowSizeCmd)
}

func (m rootModel) SetPreviousModel() (tea.Model, tea.Cmd) {
	prevModel := m.history.Pop()

	if prevModel != nil {
		m.currentModel = prevModel
		return m, tea.WindowSize()
	}

	return m, tea.Quit
}

func (m rootModel) View() string {
	return m.currentModel.View()
}

// Run initializes a new Program with rootModel and run it. Return the final model.
func Run(app ports.Application, initialModel tea.Model, logger *slog.Logger) (returnModel tea.Model, err error) {
	m := New(app, initialModel, logger)

	return tea.NewProgram(
		m, tea.WithAltScreen(),
	).Run()
}

type history struct {
	stack []tea.Model
}

func (h *history) Pop() tea.Model {
	if h.lastIndex() == -1 {
		return nil
	}
	last := h.stack[h.lastIndex()]
	h.stack = h.stack[:h.lastIndex()]

	return last
}

func (h *history) Push(m tea.Model) {
	if m != nil {
		h.stack = append(h.stack, m)
	}
}

// lastIndex return the last element of the queue index.
// If the queue is empty, -1 is return.
func (h *history) lastIndex() int {
	if h.len() == 0 {
		return -1
	}

	return h.len() - 1
}

func (h *history) len() int {
	return len(h.stack)
}
