package sqlite

import (
	"context"
	"database/sql"
	"testing"

	"github.com/aboubakary833/codexa/internal/adapters/sqlite/migrate"
	"github.com/aboubakary833/codexa/internal/domain"
	testutils "github.com/aboubakary833/codexa/utils/tests"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

// Create a temporary sqlite test database in memory
func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory?_loc=Local&parseTime=true")

	if err != nil {
		t.Fatal(err)
	}

	// In case there is already the snippets table
	migrate.Down(db)

	err = migrate.Up(db)

	if err != nil {
		t.Fatal(err)
	}

	return db
}

func TestRepositories(t *testing.T) {
	db := setupTestDB(t)
	snippetRepository := NewSnippetRepository(db)
	techRepository := NewTechRepository(db)
	ctx := context.Background()

	t.Cleanup(func() {
		err := migrate.Down(db)

		if err != nil {
			t.Fatal(err)
		}
	})

	categories := testutils.GetCategories()
	snippets := testutils.GetSnippets()

	t.Run("techRepository Store method should store tech categories", func(t *testing.T) {
		for _, category := range categories {
			err := techRepository.Store(ctx, category)
			assert.NoError(t, err)
		}
	})

	t.Run("snippetRepository Store method should store snippets", func(t *testing.T) {
		for _, snippet := range snippets {
			err := snippetRepository.Store(ctx, snippet)
			assert.NoError(t, err)
		}
	})

	t.Run("techRepository FindAll method should return all the tech categories", func(t *testing.T) {
		expectedLength := len(categories)
		actual, err := techRepository.FindAll(ctx)

		if assert.NoError(t, err) && assert.Equal(t, expectedLength, len(actual)) {
			for i := range expectedLength {
				assert.Equal(t, categories[i].ID, actual[i].ID)
				assert.Equal(t, categories[i].Name, actual[i].Name)
			}
		}
	})

	t.Run("techRepository Retrieve method should return correct tech categories", func(t *testing.T) {
		tests := []struct {
			needle       string
			expectedID   string
			expectedName string
		}{
			{
				needle:       "nodejs",
				expectedID:   "javascript",
				expectedName: "JavaScript",
			},
			{
				needle:       "js",
				expectedID:   "javascript",
				expectedName: "JavaScript",
			},
			{
				needle:       "golang",
				expectedID:   "go",
				expectedName: "Go",
			},
			{
				needle:       "php",
				expectedID:   "php",
				expectedName: "PHP",
			},
		}

		for _, test := range tests {
			category, err := techRepository.Retrieve(ctx, test.needle)

			if assert.NoError(t, err) {
				assert.Equal(t, test.expectedID, category.ID)
				assert.Equal(t, test.expectedName, category.Name)
			}
		}
	})

	t.Run("techRepository GetAliases should return JavaScript category aliases", func(t *testing.T) {
		aliases, err := techRepository.GetAliases(ctx, "javascript")

		if assert.NoError(t, err) {
			assert.Equal(t, categories[1].Aliases, aliases)
		}
	})

	t.Run("snippetRepository FindAll method should return all the snippets", func(t *testing.T) {
		expectedLength := len(snippets)
		actual, err := snippetRepository.FindAll(ctx)
		if assert.NoError(t, err) &&
			assert.Equal(t, expectedLength, len(actual)) {

			for i := range expectedLength {
				assertEntriesEqual(t, snippets[i], actual[i])
			}
		}
	})

	t.Run("snippetRepository Retrieve method should find JavaScript Arrays snippet entry", func(t *testing.T) {
		expected := snippets[2]
		snippet, err := snippetRepository.Retrieve(ctx, "javascript", "arrays")

		if assert.NoError(t, err) {
			assertEntriesEqual(t, expected, snippet)
		}
	})

	t.Run("snippetRepository Retrieve method should return domain.ErrEntryNotFound", func(t *testing.T) {
		_, err := snippetRepository.Retrieve(ctx, "javascript", "proxy")
		assert.ErrorIs(t, err, domain.ErrSnippetNotFound)
	})

	t.Run("snippetRepository FindAllByTech method should get all snippets for a given tech category", func(t *testing.T) {
		tests := []struct {
			expected []domain.Snippet
			TechID   string
		}{
			{expected: snippets[:2], TechID: "go"},
			{expected: snippets[2:4], TechID: "javascript"},
		}

		for _, test := range tests {
			actual, err := snippetRepository.FindAllByTech(ctx, test.TechID)

			if assert.NoError(t, err) && assert.Equal(t, len(test.expected), len(actual)) {
				for i := range len(test.expected) {
					assertEntriesEqual(t, test.expected[i], actual[i])
				}
			}
		}
	})

	t.Run("snippetRepository Search method should return all PHP snippets", func(t *testing.T) {
		expected := snippets[4:]
		actual, err := snippetRepository.Search(ctx, "php", "")

		if assert.NoError(t, err) {
			for i := range len(expected) {
				assertEntriesEqual(t, expected[i], actual[i])
			}
		}
	})

	t.Run("snippetRepository Search method should return Go slices snippet entry only", func(t *testing.T) {
		expected := snippets[:1]
		actual, err := snippetRepository.Search(ctx, "go", "slices")

		if assert.NoError(t, err) {
			for i := range len(expected) {
				assertEntriesEqual(t, expected[i], actual[i])
			}
		}
	})
}

func assertEntriesEqual(t *testing.T, expected, actual domain.Snippet) bool {
	assert := assert.New(t)

	return assert.Equal(expected.ID, actual.ID) &&
		assert.Equal(expected.TechID, actual.TechID) &&
		assert.Equal(expected.Topic, actual.Topic) &&
		assert.Equal(expected.Filepath, actual.Filepath) &&
		assert.True(expected.CreatedAt.Equal(actual.CreatedAt)) &&
		assert.True(expected.UpdatedAt.Equal(actual.UpdatedAt))
}

