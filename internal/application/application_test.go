package application

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/aboubakary833/codexa/internal/domain"
	"github.com/aboubakary833/codexa/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestMethodsWithSimpleCalls(t *testing.T) {
	snippetRepository := mocks.NewSnippetRepository()
	techRepository := mocks.NewTechRepository()
	registry := mocks.NewRegistry()
	fetcher := mocks.NewFetcher()

	app := New(snippetRepository, techRepository, registry, fetcher)
	ctx := context.Background()

	t.Run("ListTechCategories should return a list of tech categories", func(t *testing.T) {
		categories := []domain.Tech{
			{
				ID:   "go",
				Name: "Go",
				Aliases: []domain.TechAlias{
					{ID: "golang", TechID: "go", Name: "Golang"},
				},
			},
			{
				ID:   "javascript",
				Name: "JavaScript",
				Aliases: []domain.TechAlias{
					{ID: "js", TechID: "javascript", Name: "JS"},
					{ID: "nodejs", TechID: "javascript", Name: "NodeJS"},
				},
			},
			{
				ID:   "php",
				Name: "PHP",
				Aliases: []domain.TechAlias{
					{ID: "php8", TechID: "php", Name: "PHP8"},
				},
			},
		}
		techRepository.On("FindAll", ctx).Return(categories, nil).Once()

		res, err := app.ListTechCategories(ctx)

		if assert.NoError(t, err) {
			assert.Equal(t, categories, res)
			techRepository.AssertExpectations(t)
		}
	})

	t.Run("ListSnippets should return a given tech category snippets", func(t *testing.T) {
		entries := []domain.Snippet{
			{
				ID:       "go:slices",
				TechID:   "go",
				Topic:    "slices",
				Filepath: "go/slices.md",
			},
			{
				ID:       "go:maps",
				TechID:   "go",
				Topic:    "maps",
				Filepath: "go/maps.md",
			},
		}

		snippetRepository.On("FindAllByTech", ctx, "go").Return(entries, nil).Once()

		res, err := app.ListSnippets(ctx, "go")

		if assert.NoError(t, err) {
			assert.Equal(t, entries, res)
			snippetRepository.AssertExpectations(t)
		}
	})

	t.Run("GetSnippetContent should load go slices entry", func(t *testing.T) {
		entry := domain.Snippet{
			ID:       "go:slices",
			TechID:   "go",
			Topic:    "slices",
			Filepath: "go/slices.md",
		}

		content := "This is slices entry content"
		registry.On("LoadSnippet", ctx, entry).Return(content, nil).Once()

		res, err := app.GetSnippetContent(ctx, entry)

		if assert.NoError(t, err) {
			assert.Equal(t, content, res)
			registry.AssertExpectations(t)
		}
	})

	t.Run("FindSnippet should return ErrEntryNotFound when tech category is not found", func(t *testing.T) {
		techRepository.On("Retrieve", ctx, "Ocaml").Return(nil, domain.ErrTechNotFound).Once()

		res, err := app.FindSnippet(ctx, "Ocaml", "strings")
		if assert.Empty(t, res) {
			assert.ErrorIs(t, err, domain.ErrSnippetNotFound)
			techRepository.AssertExpectations(t)
			snippetRepository.AssertNotCalled(t, "Retrieve")
		}
	})

	t.Run("FindSnippet should return ErrEntryNotFound when snippet is not found", func(t *testing.T) {
		category := domain.Tech{
			ID:   "go",
			Name: "Go",
			Aliases: []domain.TechAlias{
				{ID: "golang", TechID: "go", Name: "Golang"},
			},
		}
		techRepository.On("Retrieve", ctx, "Go").Return(category, nil).Once()
		snippetRepository.On("Retrieve", ctx, "go", "slices").Return(nil, domain.ErrSnippetNotFound).Once()

		res, err := app.FindSnippet(ctx, "Go", "slices")

		if assert.Empty(t, res) {
			assert.ErrorIs(t, err, domain.ErrSnippetNotFound)
			snippetRepository.AssertExpectations(t)
		}
	})

	t.Run("ListRemoteTechSnippets should return a list of remote snippets", func(t *testing.T) {
		snippets := []domain.Snippet{
			{
				ID:       "go:slices",
				Topic:    "slices",
				Filepath: "go/slices.md",
			},
			{
				ID:       "go:maps",
				Topic:    "maps",
				Filepath: "go/maps.md",
			},
		}
		fetcher.On("PullTechSnippets", ctx, "go").Return(snippets, nil).Once()

		actual, err := app.ListRemoteTechSnippets(ctx, domain.RemoteTech{Dirname: "go"})

		if assert.NoError(t, err) {
			assert.Equal(t, snippets, actual)
		}
	})
}

func TestListRemoteTechCategoriesMethod(t *testing.T) {
	snippetRepository := mocks.NewSnippetRepository()
	techRepository := mocks.NewTechRepository()
	registry := mocks.NewRegistry()
	fetcher := mocks.NewFetcher()

	app := New(snippetRepository, techRepository, registry, fetcher)
	ctx := context.Background()

	manifest := domain.Manifest{
		Version: "1.0",
		Techs: []domain.RemoteTech{
			{
				Tech: domain.Tech{
					ID:   "go",
					Name: "Go",
					Aliases: []domain.TechAlias{
						{ID: "golang", Name: "Golang"},
					},
				},
				Dirname: "go",
			},
			{
				Tech: domain.Tech{
					ID:   "javascript",
					Name: "JavaScript",
					Aliases: []domain.TechAlias{
						{ID: "js", Name: "JS"},
						{ID: "nodejs", Name: "NodeJS"},
					},
				},
				Dirname: "js",
			},
		},
	}

	t.Run("it should return cached manifest categories", func(t *testing.T) {
		cachedManifest := domain.CachedManifest{
			Manifest:  manifest,
			UpdatedAt: time.Now().Add(-10 * time.Minute),
		}
		registry.On("GetManifest", ctx).Return(cachedManifest, nil).Once()

		categories, err := app.ListRemoteTechCategories(ctx)

		if assert.NoError(t, err) {
			registry.AssertExpectations(t)
			assert.Equal(t, cachedManifest.Techs, categories)
		}
	})

	t.Run("it should fetch remote registry when cache not found", func(t *testing.T) {
		registry.On("GetManifest", ctx).Return(nil, domain.ErrCachedManifestNotFound).Once()
		registry.On("CreateOrUpdateManifest", ctx, manifest).Return(nil).Once()
		fetcher.On("PullManifest", ctx).Return(manifest, nil).Once()

		categories, err := app.ListRemoteTechCategories(ctx)

		if assert.NoError(t, err) {
			fetcher.AssertExpectations(t)
			registry.AssertExpectations(t)
			assert.Equal(t, manifest.Techs, categories)
		}
	})

	t.Run("it should fetch remote registry when cache is empty", func(t *testing.T) {
		cachedManifest := domain.CachedManifest{
			Manifest: domain.Manifest{
				Version: "1.0",
				Techs:   []domain.RemoteTech{},
			},
			UpdatedAt: time.Now().Add(-10 * time.Minute),
		}
		registry.On("GetManifest", ctx).Return(cachedManifest, nil).Once()
		registry.On("CreateOrUpdateManifest", ctx, manifest).Return(nil).Once()
		fetcher.On("PullManifest", ctx).Return(manifest, nil).Once()

		categories, err := app.ListRemoteTechCategories(ctx)

		if assert.NoError(t, err) {
			fetcher.AssertExpectations(t)
			registry.AssertExpectations(t)
			assert.Equal(t, manifest.Techs, categories)
		}
	})

	t.Run("it should fetch remote registry when cache expired", func(t *testing.T) {
		cachedManifest := domain.CachedManifest{
			Manifest:  manifest,
			UpdatedAt: time.Now().Add(-31 * time.Minute),
		}
		registry.On("GetManifest", ctx).Return(cachedManifest, nil).Once()
		registry.On("CreateOrUpdateManifest", ctx, manifest).Return(nil).Once()
		fetcher.On("PullManifest", ctx).Return(manifest, nil).Once()

		categories, err := app.ListRemoteTechCategories(ctx)

		if assert.NoError(t, err) {
			fetcher.AssertExpectations(t)
			registry.AssertExpectations(t)
			assert.Equal(t, manifest.Techs, categories)
		}
	})

	t.Run("FindRemoteTechCategory should find a given tech category from categories list", func(t *testing.T) {
		cachedManifest := domain.CachedManifest{
			Manifest:  manifest,
			UpdatedAt: time.Now().Add(-10 * time.Minute),
		}
		registry.On("GetManifest", ctx).Return(cachedManifest, nil).Once()

		tech, err := app.FindRemoteTechCategory(ctx, "golang")

		if assert.NoError(t, err) {
			assert.Equal(t, manifest.Techs[0], tech)
		}
	})
}

func TestSyncsMethod(t *testing.T) {
	snippetRepository := mocks.NewSnippetRepository()
	techRepository := mocks.NewTechRepository()
	registry := mocks.NewRegistry()
	fetcher := mocks.NewFetcher()
	
	app := New(
		snippetRepository, techRepository,
		registry, fetcher,
	)
	ctx := context.Background()

	s := domain.Snippet{
		ID: "go:context",
		TechID: "go",
		Topic: "Go context",
		Filepath: "go/context.md",
	}

	t.Run("SyncSnippet should create/update a given snippet", func(t *testing.T) {
		content := "Snippet content from remote registry"

		copySnippet := s
		copySnippet.TechID = "go"

		snippetRepository.On("CreateOrUpdate", ctx, &copySnippet).Return(nil).Once()
		fetcher.On("PullSnippetContent", ctx, "go/context.md").Return(content, nil).Once()
		registry.On("CreateOrUpdateSnippet", ctx, "go/context.md", content).Return(nil).Once()

		if err := app.SyncSnippet(ctx, s); assert.NoError(t, err) {
			snippetRepository.AssertExpectations(t)
			registry.AssertExpectations(t)
			fetcher.AssertExpectations(t)
		}
	})

	t.Run("SyncSnippet should rollback when registry fail creating/updating snippet content", func(t *testing.T) {
		content := "Snippet content from remote registry"

		copySnippet := s
		copySnippet.TechID = "go"

		snippetRepository.On("Delete", ctx, copySnippet).Return(nil).Once()
		snippetRepository.On("CreateOrUpdate", ctx, &copySnippet).Return(nil).Once()
		fetcher.On("PullSnippetContent", ctx, "go/context.md").Return(content, nil).Once()
		registry.On("CreateOrUpdateSnippet", ctx, "go/context.md", content).Return(os.ErrPermission).Once()

		err := app.SyncSnippet(ctx, s)

		assert.ErrorIs(t, err, os.ErrPermission)
		snippetRepository.AssertExpectations(t)
		registry.AssertExpectations(t)
		fetcher.AssertExpectations(t)

	})

}

func TestParseSearchInputMethod(t *testing.T) {
	app := &app{}

	tests := []struct {
		input            string
		expectedCategory string
		expectedTopic    string
	}{
		{" go/", "go", ""},
		{" / go slice ", "go", "slice"},
		{"/js array", "js", "array"},
		{"php/ namespace ", "php", "namespace"},
		{"php/class", "php", "class"},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("ParseSearchInput test #%d", i+1), func(t *testing.T) {
			category, topic := app.ParseSearchInput(test.input)

			assert.Equal(t, test.expectedCategory, category)
			assert.Equal(t, test.expectedTopic, topic)
		})
	}
}
