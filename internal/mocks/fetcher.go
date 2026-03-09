package mocks

import (
	"context"

	"github.com/aboubakary833/codexa/internal/domain"
	"github.com/stretchr/testify/mock"
)

type fetcher struct {
	mock.Mock
}

func NewFetcher() *fetcher {
	return new(fetcher)
}

func (f *fetcher) PullManifest(ctx context.Context) (domain.Manifest, error) {
	args := f.Called(ctx)

	if args.Get(0) == nil {
		return domain.Manifest{}, args.Error(1)
	}

	return args.Get(0).(domain.Manifest), nil
}

func (f *fetcher) PullTechSnippets(ctx context.Context, tech string) ([]domain.Snippet, error) {
	args := f.Called(ctx, tech)

	if args.Get(0) == nil {
		return []domain.Snippet{}, args.Error(1)
	}

	return args.Get(0).([]domain.Snippet), nil
}

func (f *fetcher) PullSnippetContent(ctx context.Context, snippetPath string) (string, error) {
	args := f.Called(ctx, snippetPath)

	if args.Get(0) == nil {
		return "", args.Error(1)
	}

	return args.String(0), nil
}
