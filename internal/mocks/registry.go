package mocks

import (
	"context"

	"github.com/aboubakary833/codexa/internal/domain"
	"github.com/stretchr/testify/mock"
)

type Registry struct {
	mock.Mock
}

func NewRegistry() *Registry {
	return new(Registry)
}

func (registry *Registry) LoadSnippet(ctx context.Context, snippet domain.Snippet) (string, error) {
	args := registry.Called(ctx, snippet)

	return args.String(0), args.Error(1)
}

func (registry *Registry) CreateOrUpdateSnippet(ctx context.Context, path, content string) error {
	args := registry.Called(ctx, path, content)
	return args.Error(0)
}

func (registry *Registry) GetManifest(ctx context.Context) (domain.CachedManifest, error) {
	args := registry.Called(ctx)

	if args.Get(0) == nil {
		return domain.CachedManifest{}, args.Error(1)
	}

	return args.Get(0).(domain.CachedManifest), nil
}

func (registry *Registry) CreateOrUpdateManifest(ctx context.Context, manifest domain.Manifest) error {
	args := registry.Called(ctx, manifest)
	return args.Error(0)
}

func (registry *Registry) Stat() error {
	args := registry.Called()
	return args.Error(0)
}
