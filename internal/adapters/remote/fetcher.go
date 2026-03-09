package remote

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"github.com/aboubakary833/codexa/internal/domain"
	"github.com/aboubakary833/codexa/utils"
)

type remoteSnippet struct {
	ID       string `json:"id"`
	Topic    string `json:"topic"`
	Filename string `json:"filename"`
}

type fetcher struct {
	url    string
	client *http.Client
}

func NewFetcher(remoteUrl string) *fetcher {
	return &fetcher{
		url: remoteUrl,
		client: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

// PullManifest get remote registry manifest
func (f *fetcher) PullManifest(ctx context.Context) (domain.Manifest, error) {
	res, err := f.doHttpGetRequest(ctx, "/registry.json", domain.ErrManifestNotFound)

	if err != nil {
		return domain.Manifest{}, err
	}
	defer res.Body.Close()

	var manifest domain.Manifest
	decoder := json.NewDecoder(res.Body)

	if err = decoder.Decode(&manifest); err != nil {
		return domain.Manifest{}, err
	}

	return manifest, nil
}

func (f *fetcher) PullTechSnippets(ctx context.Context, dirname string) ([]domain.Snippet, error) {
	path := fmt.Sprintf("/%s/index.json", dirname)
	res, err := f.doHttpGetRequest(ctx, path, domain.ErrTechNotFound)

	if err != nil {
		return []domain.Snippet{}, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var rs []remoteSnippet

	if err = decoder.Decode(&rs); err != nil {
		return []domain.Snippet{}, err
	}

	snippets := utils.Mutate(rs, func(i remoteSnippet) domain.Snippet {
		path := filepath.Join(dirname, i.Filename)
		return domain.Snippet{
			ID: i.ID,
			Topic: i.Topic,
			Filepath: path,
		}
	})

	return snippets, nil
}

// PullSnippetContent pull a specific snippet content from the remote registry
func (f *fetcher) PullSnippetContent(ctx context.Context, snippetUrl string) (string, error) {
	res, err := f.doHttpGetRequest(ctx, snippetUrl, domain.ErrSnippetNotFound)

	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	rawBytes, err := io.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	return string(rawBytes), nil
}

// doHttpGetRequest is a util function for fetcher http GET requests
func (f *fetcher) doHttpGetRequest(ctx context.Context, path string, notFoundError error) (*http.Response, error) {
	url, err := url.JoinPath(f.url, path)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "codexa-cli/1.0")
	req.Header.Set("Accept", "application/vnd.github.raw")
	res, err := f.client.Do(req)

	if err != nil {
		var netErr net.Error

		if errors.As(err, &netErr) && netErr.Timeout() {
			return nil, fmt.Errorf("request timed out")
		}

		return nil, err
	}

	switch res.StatusCode {
	case http.StatusNotFound:
		return nil, notFoundError

	case http.StatusBadRequest:
		return nil, fmt.Errorf("invalid request")
	}

	if res.StatusCode > 400 {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return res, nil
}
