package cli

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/aboubakary833/codexa/internal/adapters/tui"
	"github.com/aboubakary833/codexa/internal/adapters/tui/models"
	"github.com/aboubakary833/codexa/internal/domain"
	"github.com/aboubakary833/codexa/internal/ports"
	"github.com/aboubakary833/codexa/utils"
	"github.com/briandowns/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	timeoutError  = errors.New("request timed out. Please retry.")
	internalError = errors.New("an unexpected error occurred.")
)

type controller struct {
	app    ports.Application
	logger *slog.Logger
}

// render launch the TUI with an initial view
func (c controller) render(initialModel tea.Model) error {
	_, err := tui.Run(c.app, initialModel, c.logger)
	return err
}

// renderSnippetContent retrieve a given entry, loads its content and render it
func (c controller) renderSnippetContent(tech, topic string) error {
	entry, err := c.findSnippet(tech, topic)

	if err != nil {
		return err
	}

	content, err := c.getSnippetContent(entry)

	if err != nil {
		return err
	}
	m := models.NewSnippetViewModel(entry.TechID, entry.Topic, content)
	return c.render(m)
}

// renderTechSnippetsList find a list of entries related to a given category and render it
func (c controller) renderTechSnippetsList(categoryName string) error {
	category, err := c.findTech(categoryName)

	if err != nil {
		return err
	}

	snippets, err := c.getTechSnippets(category)

	if err != nil {
		return err
	}

	m := models.NewEntries(category, snippets)
	return c.render(m)
}

func (c controller) syncSnippet(tech, snippetName string) error {
	s := spinner.New(spinner.CharSets[14], 120*time.Millisecond)
	s.Suffix = " Checking registry..."
	s.Start()
	rt, err := c.findRemoteTech(tech)

	if err != nil {
		s.Stop()
		return err
	}

	s.Suffix = fmt.Sprintf("Searching for %s...", snippetName)
	s.Restart()

	rs, err := c.findRemoteSnippet(rt, snippetName)

	if err != nil {
		s.Stop()
		return err
	}

	snippetPath := fmt.Sprintf("%s/%s", strings.ToLower(rt.Name), strings.ToLower(rs.Topic))
	s.Suffix = "Pulling " + snippetPath + "..."
	s.Restart()

	err = c.download(rt, rs)
	s.Stop()

	if err != nil {
		return err
	}

	fmt.Printf("%s synced", snippetPath)
	return nil
}

// syncTechCategory sync a given category snippets
func (c controller) syncTechCategory(tech string) error {
	s := spinner.New(spinner.CharSets[14], 120*time.Millisecond)
	s.Suffix = " Checking registry..."
	s.Start()

	rt, err := c.findRemoteTech(tech)

	if err != nil {
		s.Stop()
		return err
	}

	s.Suffix = fmt.Sprintf("Syncing %s...", rt.ID)
	s.Restart()

	rs, err := c.listRemoteTechSnippets(rt)

	if err != nil || len(rs) == 0 {
		s.Stop()
		return err
	}

	s.Suffix = fmt.Sprintf("Pulling %s snippets(%d)...", rt.ID, len(rs))
	s.Restart()

	err = c.download(rt, rs...)
	s.Stop()

	if err != nil {
		return err
	}

	fmt.Printf("✔ %d %s snippets synced\n", len(rs), rt.ID)
	return nil
}

// download Download one or more snippet for given tech category
func (c controller) download(rt domain.RemoteTech, snippets ...domain.Snippet) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	return c.app.Sync(ctx, rt, snippets...)
}

func (c controller) findRemoteSnippet(rt domain.RemoteTech, needle string) (domain.Snippet, error) {

	snippets, err := c.listRemoteTechSnippets(rt)

	if err != nil {
		return domain.Snippet{}, err
	}

	for _, s := range snippets {
		if s.Match(needle) {
			return s, nil
		}
	}
	err = fmt.Errorf("%s %q %s", rt.ID, needle, domain.ErrSnippetNotFound.Error())

	return domain.Snippet{}, err
}

func (c *controller) listRemoteTechSnippets(rt domain.RemoteTech) ([]domain.Snippet, error) {
	return utils.Exec(func(ctx context.Context) ([]domain.Snippet, error) {
		return c.app.ListRemoteTechSnippets(ctx, rt)
	})
}

// findRemoteTech get the tech category from the cache or remote registry
func (c controller) findRemoteTech(categoryName string) (domain.RemoteTech, error) {
	rt, err := utils.Exec(func(ctx context.Context) (domain.RemoteTech, error) {
		return c.app.FindRemoteTechCategory(ctx, categoryName)
	})

	if err != nil {
		if errors.Is(err, domain.ErrTechNotFound) {
			return domain.RemoteTech{}, fmt.Errorf("%q %s in the registry", categoryName, err.Error())
		}

		return domain.RemoteTech{}, err
	}

	return rt, nil
}

// getCategorySnippets is a helper function that return a list of entries for a given category
func (c controller) getTechSnippets(tech domain.Tech) ([]domain.Snippet, error) {
	snippets, err := utils.Exec(func(ctx context.Context) ([]domain.Snippet, error) {
		return c.app.ListSnippets(ctx, tech.ID)
	})

	if err != nil {
		return []domain.Snippet{}, c.handleUnexpectedError(err)
	}

	return snippets, nil
}

// findTech is a helper function that find a tech category and return it
func (c controller) findTech(techName string) (domain.Tech, error) {
	category, err := utils.Exec(func(ctx context.Context) (domain.Tech, error) {
		return c.app.FindTechCategory(ctx, techName)
	})

	if err != nil {
		if errors.Is(err, domain.ErrTechNotFound) {
			errStr := err.Error()
			return domain.Tech{}, fmt.Errorf("%q snippets %s", techName, errStr)
		}

		return domain.Tech{}, c.handleUnexpectedError(err)
	}

	return category, nil
}

// findSnippet is a helper function that find and return a snippet
func (c controller) findSnippet(tech, topic string) (domain.Snippet, error) {

	snippet, err := utils.Exec(func(ctx context.Context) (domain.Snippet, error) {
		return c.app.FindSnippet(ctx, tech, topic)
	})

	if err != nil {
		if errors.Is(err, domain.ErrSnippetNotFound) {
			title := fmt.Sprintf("%s %s", tech, topic)

			return domain.Snippet{}, fmt.Errorf("%q %s", title, err.Error())
		}

		return domain.Snippet{}, c.handleUnexpectedError(err)
	}

	return snippet, nil
}

// getSnippetContent is a helper function that load and return a given snippet content
func (c controller) getSnippetContent(snippet domain.Snippet) (string, error) {
	content, err := utils.Exec(func(ctx context.Context) (string, error) {
		return c.app.GetSnippetContent(ctx, snippet)
	})

	if err != nil {
		if errors.Is(err, domain.ErrSnippetContentNotFound) ||
			errors.Is(err, domain.ErrSnippetContentCantBeRead) {
			return "", err
		}

		return "", c.handleUnexpectedError(err)
	}

	return content, nil
}

// handleUnexpectedError is a helper function for handling undesired errors
func (c controller) handleUnexpectedError(err error) error {
	if errors.Is(err, context.DeadlineExceeded) {
		return timeoutError
	}

	c.logger.Error(err.Error(), slog.Any("Error", err))

	return internalError
}
