package application

import (
	"context"
	"errors"
	"strings"

	"github.com/aboubakary833/codexa/internal/domain"
	"github.com/aboubakary833/codexa/internal/ports"
)

type app struct {
	snippetRepository ports.SnippetRepository
	techRepository    ports.TechRepository
	registry          ports.Registry
	fetcher			  ports.Fetcher
}

func New(
	snippetRepository ports.SnippetRepository,
	techRepository ports.TechRepository,
	registry ports.Registry,
) *app {
	return &app{
		snippetRepository: snippetRepository,
		techRepository:    techRepository,
		registry:          registry,
	}
}

// ListTechCategories return all the available tech categories.
func (app *app) ListTechCategories(ctx context.Context) ([]domain.Tech, error) {
	return app.techRepository.FindAll(ctx)
}

// ListSnippets return all entries for a specific tech category
func (app *app) ListSnippets(ctx context.Context, techID string) ([]domain.Snippet, error) {
	if techID == "" {
		return []domain.Snippet{}, domain.ErrTechNotFound
	}

	return app.snippetRepository.FindAllByTech(ctx, techID)
}

// FindTechCategory try to retrieve a tech category
func (app *app) FindTechCategory(ctx context.Context, tech string) (domain.Tech, error) {
	if tech == "" {
		return domain.Tech{}, domain.ErrTechNotFound
	}

	return app.techRepository.Retrieve(ctx, tech)
}

// FindEntry retrieve a single snippet from the repository by first trying to find the tech category
func (app *app) FindSnippet(ctx context.Context, tech, topic string) (domain.Snippet, error) {
	if topic == "" {
		return domain.Snippet{}, domain.ErrSnippetNotFound
	}
	category, err := app.FindTechCategory(ctx, tech)

	if err != nil {

		if errors.Is(err, domain.ErrTechNotFound) {
			return domain.Snippet{}, domain.ErrSnippetNotFound
		}

		return domain.Snippet{}, err
	}

	return app.snippetRepository.Retrieve(ctx, category.ID, topic)
}

// GetSnippetContent retrieve and load a snippet content and return it.
func (app *app) GetSnippetContent(ctx context.Context, snippet domain.Snippet) (string, error) {
	return app.registry.LoadContent(ctx, snippet)
}

// Search take a non empty input string, make a search to the repository and return a slice of Snippets
func (app *app) Search(ctx context.Context, input string) ([]domain.Snippet, error) {
	tech, topic := app.ParseSearchInput(input)

	if tech == "" {
		return app.snippetRepository.FindAll(ctx)
	}

	return app.snippetRepository.Search(ctx, tech, topic)
}

func (app *app) SyncSnippet(ctx context.Context, snippet domain.Snippet) error {
	
	return nil
}

// ParseSearchInput parse a provided input string into tech category and topic
func (app *app) ParseSearchInput(input string) (tech, topic string) {
	input = strings.Trim(input, " /")

	if input == "" {
		return
	}

	if strings.Contains(input, " ") {
		parts := strings.Fields(input)
		if len(parts) >= 1 {
			tech = strings.Trim(parts[0], "/")
		}
		if len(parts) >= 2 {
			topic = strings.Trim(parts[1], "/")
		}
		return
	}

	parts := strings.Split(input, "/")
	if len(parts) >= 1 {
		tech = parts[0]
	}
	if len(parts) >= 2 {
		topic = parts[1]
	}

	return
}
