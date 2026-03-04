package tui

import (
	"context"
	"errors"

	"github.com/aboubakary833/codexa/internal/adapters/tui/models"
	"github.com/aboubakary833/codexa/internal/domain"
	"github.com/aboubakary833/codexa/internal/ports"
	"github.com/aboubakary833/codexa/utils"
	tea "github.com/charmbracelet/bubbletea"
)

type controller struct {
	app ports.Application
}

// newController return a new controller for the rootModel
func newController(application ports.Application) controller {
	return controller{
		app: application,
	}
}

// loadTechCategories get a list of all available tech categories
func (c controller) loadTechCategories() tea.Cmd {

	categories, err := utils.Exec(func(ctx context.Context) ([]domain.Tech, error) {
		return c.app.ListTechCategories(ctx)
	})

	if err != nil {
		return c.sendErrorCmd(err)
	}

	return func() tea.Msg {
		return models.TechsLoadedMsg{
			Techs: categories,
		}
	}
}

// loadTechSnippets get list of snippets for a specific tech category
func (c controller) loadTechSnippets(tech domain.Tech) tea.Cmd {

	snippet, err := utils.Exec(func(ctx context.Context) ([]domain.Snippet, error) {
		return c.app.ListSnippets(ctx, tech.ID)
	})

	if err != nil {
		return c.sendErrorCmd(err)
	}

	return func() tea.Msg {
		return models.TechSnippetsLoadedMsg{
			Tech:     tech,
			Snippets: snippet,
		}
	}
}

// loadSnippetContent get a specific snippet content from the application
func (c controller) loadSnippetContent(snippet domain.Snippet) tea.Cmd {

	content, err := utils.Exec(func(ctx context.Context) (string, error) {
		return c.app.GetSnippetContent(ctx, snippet)
	})

	if err != nil {
		return c.sendErrorCmd(err)
	}

	return func() tea.Msg {
		return models.SnippetLoadedMsg{
			Tech:    snippet.TechID,
			Topic:   snippet.Topic,
			Content: content,
		}
	}
}

// search make a search request to the application
func (c controller) search(query string) tea.Cmd {

	result, err := utils.Exec(func(ctx context.Context) ([]domain.Snippet, error) {
		return c.app.Search(ctx, query)
	})

	if err != nil {
		return c.sendErrorCmd(err)
	}

	return func() tea.Msg {
		return models.SearchResponseMsg{
			Snippets: result,
		}
	}
}

// sendErrorCmd return a cmd that signal the system
// when expected or unexpected error occur
func (c controller) sendErrorCmd(err error) tea.Cmd {

	if errors.Is(err, context.DeadlineExceeded) {
		return func() tea.Msg {
			return TimeoutErrorMsg{}
		}
	}

	//NOTE: Log unknown error message to stderr

	return func() tea.Msg {
		return InternalErrorMsg{}
	}
}
