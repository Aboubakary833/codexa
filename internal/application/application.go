package application

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/aboubakary833/codexa/internal/domain"
	"github.com/aboubakary833/codexa/internal/ports"
	"golang.org/x/sync/errgroup"
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
	return app.registry.LoadSnippet(ctx, snippet)
}

// Search take a non empty input string, make a search to the repository and return a slice of Snippets
func (app *app) Search(ctx context.Context, input string) ([]domain.Snippet, error) {
	tech, topic := app.ParseSearchInput(input)

	if tech == "" {
		return app.snippetRepository.FindAll(ctx)
	}

	return app.snippetRepository.Search(ctx, tech, topic)
}

// Download create or update one or more snippets. The function also create
// the tech category if it does'nt exists
func (app *app) Download(ctx context.Context, rt domain.RemoteTech, snippets ...domain.RemoteSnippet) error {
	_, err := app.techRepository.FindByID(ctx, rt.ID)

	if err != nil && !errors.Is(err, domain.ErrTechNotFound) {
		return err
	}

	tech := rt.Tech

	if err != nil {
		err = app.techRepository.Store(ctx, tech)

		if err != nil {
			return err
		}
	}
	g, gCtx := errgroup.WithContext(ctx)
	g.SetLimit(5)

	for _, rs := range snippets {
		g.Go(func () error {
			return app.downloadSnippet(gCtx, rt, rs)
		})
	}

	return g.Wait()
}

// downloadSnippet is a helper function to Download syncing a snippet
func (app *app) downloadSnippet(ctx context.Context, rt  domain.RemoteTech, rs domain.RemoteSnippet) error {
	remotePath := fmt.Sprintf("%s/%s", rt.Dirname, rs.Filename)
	content, err := app.fetcher.PullSnippetContent(ctx, remotePath)

	if err != nil {
		return err
	}

	filePath := filepath.Join(rt.Dirname, rs.Filename)
	snippet := &domain.Snippet{
		ID: rs.ID,
		TechID: rt.ID,
		Topic: rs.Topic,
		Filepath: filePath,
	}

	if err = app.snippetRepository.CreateOrUpdate(ctx, snippet); err != nil {
		return err
	}

	if err = app.registry.CreateOrUpdateSnippet(ctx, filePath, content); err != nil {
		app.snippetRepository.Delete(ctx, *snippet)
		return err
	}

	return nil
}

// FindRemoteTechCategory retrieve a remote tech category
func (app *app) FindRemoteTechCategory(ctx context.Context, tech string) (domain.RemoteTech, error) {
	categories, err := app.ListRemoteTechCategories(ctx)

	if err != nil {
		return domain.RemoteTech{}, err
	}

	for _, category := range categories {
		if category.Match(tech) {
			return category, nil
		}
	}

	return domain.RemoteTech{}, domain.ErrRemoteTechNotFound
}


// FindRemoteSnippet retrieve a single remote snippet
func (app *app) FindRemoteSnippet(ctx context.Context, tech, topic string) (domain.RemoteSnippet, error) {
	rt, err := app.FindRemoteTechCategory(ctx, tech)

	if err != nil {
		if errors.Is(err, domain.ErrRemoteTechNotFound) {
			return domain.RemoteSnippet{}, domain.ErrRemoteSnippetNotFound
		}

		return domain.RemoteSnippet{}, err
	}

	snippets, err := app.ListRemoteTechSnippets(ctx, rt)

	if err != nil {
		return domain.RemoteSnippet{}, err
	}

	for _, snippet := range snippets {
		if snippet.Match(topic) {
			return snippet, nil
		}
	}

	return domain.RemoteSnippet{}, domain.ErrRemoteSnippetNotFound
}


// ListRemoteTechSnippets fetch and return a given remote tech category snippets
func (app *app) ListRemoteTechSnippets(ctx context.Context, rt domain.RemoteTech) ([]domain.RemoteSnippet, error) {
	snippets, err := app.fetcher.PullTechSnippets(ctx, rt.Dirname)

	if err != nil {
		return []domain.RemoteSnippet{}, err
	}

	return snippets, nil
}

// ListRemoteTechCategories fetch and return remote tech categories
func (app *app) ListRemoteTechCategories(ctx context.Context) ([]domain.RemoteTech, error) {
	var manifest domain.Manifest 
	cachedManifest, err := app.registry.GetManifest(ctx)

	if err != nil && !errors.Is(err, domain.ErrCachedManifestNotFound) {
		return []domain.RemoteTech{}, err
	}

	if err != nil || !cachedManifest.IsTrustWorthy() {

		// Fetch manifest in case cached one is invalid
		manifest, err = app.fetcher.PullManifest(ctx)

		if err != nil {
			return []domain.RemoteTech{}, err
		}
		
		cachedManifest = domain.CachedManifest{
			Manifest: manifest,
			UpdatedAt: time.Now(),
		}

		app.registry.CreateOrUpdateManifest(ctx, cachedManifest)

	} else {
		manifest = cachedManifest.Manifest
	}

	if err != nil {
		return []domain.RemoteTech{}, err
	}

	return manifest.Techs, nil
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
