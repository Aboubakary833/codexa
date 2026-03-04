package remote

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/aboubakary833/codexa/internal/domain"
)

type fetcher struct {
	url    string
	client *http.Client
}

func NewFetcher(remoteUrl string) fetcher {
	return fetcher{
		url: remoteUrl,
		client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

// PullManifest get remote registry manifest
func (f fetcher) PullManifest(ctx context.Context) (domain.Manifest, error) {
	res, err := f.doHttpGetRequest(ctx, "/registry.json")

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

// PullContent pull a specific entry content from the remote registry
func (f fetcher) PullContent(ctx context.Context, entryPath string) (string, error) {
	res, err := f.doHttpGetRequest(ctx, entryPath)

	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var rawBytes []byte
	_, err = res.Body.Read(rawBytes)

	if err != nil {
		return "", err
	}

	return string(rawBytes), nil
}

// doHttpGetRequest is a util function for fetcher http GET requests
func (f fetcher) doHttpGetRequest(ctx context.Context, path string) (*http.Response, error) {
	url, err := url.JoinPath(f.url, path)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "codexa/1.0")
	return f.client.Do(req)
}
