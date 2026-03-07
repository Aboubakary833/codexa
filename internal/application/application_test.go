package application

import (
	"context"
	"fmt"
	"testing"

	"github.com/aboubakary833/codexa/internal/domain"
	"github.com/aboubakary833/codexa/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestMethodsWithRepositoriesCall(t *testing.T) {
	snippetRepository := mocks.NewSnippetRepository()
	techRepository := mocks.NewTechRepository()
	registry := mocks.NewRegistry()

	app := New(snippetRepository, techRepository, registry)
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

}

func TestParseSearchInputMethod(t *testing.T) {
	app := New(nil, nil, nil)

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
