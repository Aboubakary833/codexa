package ports

import (
	"context"

	"github.com/aboubakary833/codexa/internal/domain"
)

type Application interface {
	ListTechCategories(context.Context) ([]domain.Tech, error)
	FindTechCategory(context.Context, string) (domain.Tech, error)
	FindSnippet(ctx context.Context, tech, topic string) (domain.Snippet, error)
	ListSnippets(ctx context.Context, techID string) ([]domain.Snippet, error)
	GetSnippetContent(context.Context, domain.Snippet) (string, error)
	Search(ctx context.Context, input string) ([]domain.Snippet, error)
}

type TechRepository interface {
	// Store store a tech category
	Store(context.Context, domain.Tech) error

	// FindAll return all tech categories
	FindAll(context.Context) ([]domain.Tech, error)

	// Retrieve a tech category by ID, name or alias
	Retrieve(context.Context, string) (domain.Tech, error)

	// GetAliases return a given tech category aliases
	GetAliases(context.Context, string) ([]domain.TechAlias, error)
}

type SnippetRepository interface {
	// Store store a snippet entry
	Store(context.Context, domain.Snippet) error

	// FindAll return all the available snippet entries
	FindAll(context.Context) ([]domain.Snippet, error)

	// FindAllByTech retrieve all snippet entries for a specific tech category
	FindAllByTech(ctx context.Context, techID string) ([]domain.Snippet, error)

	// Retrieve search and return a single snippet
	Retrieve(ctx context.Context, techID, topic string) (domain.Snippet, error)

	// Search query and return entries that match the given tech category and topic
	Search(ctx context.Context, tech, topic string) ([]domain.Snippet, error)
}

type Registry interface {
	// LoadContent load and return a snippet content
	LoadContent(context.Context, domain.Snippet) (string, error)
}

type Fetcher interface {
	// PullManifest pull the tech manifest from the remote registry
	PullManifest(context.Context) (domain.Manifest, error)

	// PullTechSnippets pull a given tech category snippets list from the remote repository
	PullTechSnippets(context.Context, string) ([]domain.RemoteSnippet, error)

	// PullSnippetContent pull a specific snippet content from the remote registry
	PullSnippetContent(ctx context.Context, snippetPath string) (string, error)
}
