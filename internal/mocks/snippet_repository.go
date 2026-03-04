package mocks

import (
	"context"

	"github.com/aboubakary833/codexa/internal/domain"
	"github.com/stretchr/testify/mock"
)

type snippetRepository struct {
	mock.Mock
}

func NewSnippetRepository() *snippetRepository {
	return new(snippetRepository)
}

func (repo *snippetRepository) Store(ctx context.Context, snippet domain.Snippet) error {
	args := repo.Called(ctx, snippet)
	return args.Error(0)
}

func (repo *snippetRepository) FindAll(ctx context.Context) ([]domain.Snippet, error) {
	args := repo.Called(ctx)
	return args.Get(0).([]domain.Snippet), args.Error(1)
}

func (repo *snippetRepository) Retrieve(ctx context.Context, techID, topic string) (domain.Snippet, error) {
	args := repo.Called(ctx, techID, topic)

	if args.Get(0) == nil {
		return domain.Snippet{}, args.Error(1)
	}

	return args.Get(0).(domain.Snippet), nil
}

func (repo *snippetRepository) FindAllByTech(ctx context.Context, techID string) ([]domain.Snippet, error) {
	args := repo.Called(ctx, techID)
	return args.Get(0).([]domain.Snippet), args.Error(1)
}

func (repo *snippetRepository) Search(ctx context.Context, tech, topic string) ([]domain.Snippet, error) {
	args := repo.Called(ctx, tech, topic)
	return args.Get(0).([]domain.Snippet), args.Error(1)
}
