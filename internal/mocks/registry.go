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

func (registry *Registry) LoadContent(ctx context.Context, snippet domain.Snippet) (string, error) {
	args := registry.Called(ctx, snippet)

	return args.String(0), args.Error(1)
}
