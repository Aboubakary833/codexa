package remote

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aboubakary833/codexa/internal/domain"
	testutils "github.com/aboubakary833/codexa/utils/tests"
	"github.com/stretchr/testify/assert"
)

func TestFetcher(t *testing.T) {
	ctx := context.Background()
	server := newRemoteTestServer()
	fetcher := NewFetcher(server.URL)

	defer server.Close()

	t.Run("PullManifest should successfully fetch the manifest from the remote repo", func(t *testing.T) {
		expected := testutils.GetRemoteCategories()
		manifest, err := fetcher.PullManifest(ctx)

		if assert.NoError(t, err) {
			assert.Equal(t, "1.0", manifest.Version)
			assert.Equal(t, expected, manifest.Techs)
		}
	})

	t.Run("PullTechSnippets should successfully fetch snippets for a specific tech category", func(t *testing.T) {
		expected := testutils.GetTechRemoteSnippets("go")
		actual, err := fetcher.PullTechSnippets(ctx, "go")

		if assert.NoError(t, err) {
			assert.Equal(t, expected, actual)
		}
	})

	t.Run("PullTechSnippets should return 404 http status error", func(t *testing.T) {
		args := []string{"", "ocaml", "c", "csharp"}

		for _, arg := range args {
			_, err := fetcher.PullTechSnippets(ctx, arg)
			assert.ErrorIs(t, err, domain.ErrRemoteTechNotFound)
		}
	})

	t.Run("PullSnippetContent should successfully fetch snippet content", func(t *testing.T) {
		tests := []struct{
			url string
			resultContent string
		}{
			{"/go/slices.md", "Go slices snippet content"},
			{"/php/classes.md", "PHP classes snippet content"},
			{"/javascript/objects.md", "JavaScript objects snippet content"},
		}

		for _, test := range tests {
			content, err := fetcher.PullSnippetContent(ctx, test.url)
			if assert.NoError(t, err) {
				assert.Equal(t, test.resultContent, content)
			}
		}
	})
}

func newRemoteTestServer() *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /registry.json", func(w http.ResponseWriter, r *http.Request) {
		techCategories := testutils.GetRemoteCategories()
		manifest := domain.Manifest{
			Version:    "1.0",
			Techs: techCategories,
		}

		writeJson(w, manifest)
	})

	mux.HandleFunc("GET /{tech}/index.json", func(w http.ResponseWriter, r *http.Request) {
		tech := r.PathValue("tech")

		if tech == "" {
			http.NotFound(w, r)
			return
		}

		snippets := testutils.GetTechRemoteSnippets(tech)

		if len(snippets) == 0 {
			http.NotFound(w, r)
			return
		}

		writeJson(w, snippets)
	})

	mux.HandleFunc("GET /{tech}/{snippet}", func(w http.ResponseWriter, r *http.Request) {
		registry := map[string]string{
			"go/slices.md": "Go slices snippet content",
			"php/classes.md": "PHP classes snippet content",
			"javascript/objects.md": "JavaScript objects snippet content",
		}

		tech := r.PathValue("tech")
		snippet := r.PathValue("snippet")

		if tech == "" || snippet == "" {
			http.NotFound(w, r)
			return
		}
		key := tech + "/" + snippet
		content, ok := registry[key]

		if !ok {
			http.NotFound(w, r)
			return
		}

		w.Write([]byte(content))
	})

	return httptest.NewServer(mux)
}



// writeJson is a helper function that encode a value "v" into the response body
func writeJson(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	if err := encoder.Encode(v); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `
			{
				"message": "Internal error",
				"error": "%s"
			}
			`, err.Error(),
		)
	}
}
