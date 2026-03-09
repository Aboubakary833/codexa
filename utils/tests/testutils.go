package testutils

import (
	"slices"
	"strings"
	"time"

	"github.com/aboubakary833/codexa/internal/domain"
)

var techCategories = []domain.Tech{
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
		ID:      "php",
		Name:    "PHP",
		Aliases: []domain.TechAlias{},
	},
}

var snippets = []domain.Snippet{
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
	{
		ID:       "javascript:arrays",
		TechID:   "javascript",
		Topic:    "Arrays",
		Filepath: "js/arrays.md",
	},
	{
		ID:       "javascript:objects",
		TechID:   "javascript",
		Topic:    "Objects",
		Filepath: "js/objects.md",
	},
	{
		ID:       "php:classes",
		TechID:   "php",
		Topic:    "Classes",
		Filepath: "php/classes.md",
	},
	{
		ID:       "php:enums",
		TechID:   "php",
		Topic:    "Enums",
		Filepath: "php/enums.md",
	},
}

// GetCategories return a list of tech categories for testing purpose
func GetCategories() []domain.Tech {
	return techCategories
}

// GetSnippets return a list of snippets for testing purpose
func GetSnippets() []domain.Snippet {
	s := slices.Clone(snippets)
	now := time.Now()

	for i, snippet := range s {
		snippet.CreatedAt = now
		snippet.UpdatedAt = now

		s[i] = snippet
	}

	return s
}

// GetTechSnippets return a given category of snippets
func GetTechSnippets(techID string) []domain.Snippet {
	var (
		allSnippets = GetSnippets()
		snippets    []domain.Snippet
	)

	for _, s := range allSnippets {
		if strings.HasPrefix(s.ID, strings.ToLower(techID)) {
			snippets = append(snippets, s)
		}
	}

	return snippets
}

// GetRemoteCategories return all test remote tech categories
func GetRemoteCategories() []domain.RemoteTech {
	var rts []domain.RemoteTech

	for _, tech := range techCategories {
		rts = append(rts, domain.RemoteTech{
			Tech: tech,
			Dirname: tech.ID,
		})
	}

	return rts
}

