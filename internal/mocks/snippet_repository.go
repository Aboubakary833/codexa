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

func (repo *snippetRepository) FindByID(ctx context.Context, ID string) (domain.Snippet, error) {
	args := repo.Called(ctx, ID)

	if args.Get(0) == nil {
		return domain.Snippet{}, args.Error(1)
	}

	return args.Get(0).(domain.Snippet), nil
}

func (repo *snippetRepository) FindAll(ctx context.Context) ([]domain.Snippet, error) {
	args := repo.Called(ctx)

	if args.Get(0) == nil {
		return []domain.Snippet{}, args.Error(1)
	}

	return args.Get(0).([]domain.Snippet), nil
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
	
	if args.Get(0) == nil {
		return []domain.Snippet{}, args.Error(1)
	}

	return args.Get(0).([]domain.Snippet), nil
}

func (repo *snippetRepository) Search(ctx context.Context, tech, topic string) ([]domain.Snippet, error) {
	args := repo.Called(ctx, tech, topic)
	
	if args.Get(0) == nil {
		return []domain.Snippet{}, args.Error(1)
	}

	return args.Get(0).([]domain.Snippet), nil
}

func (repo *snippetRepository) CreateOrUpdate(ctx context.Context, snippet *domain.Snippet) error {
	args := repo.Called(ctx, snippet)
	return args.Error(0)
}

func (repo *snippetRepository) Update(ctx context.Context, snippet *domain.Snippet) error {
	args := repo.Called(ctx, snippet)
	return args.Error(0)
}

func (repo *snippetRepository) Delete(ctx context.Context, snippet domain.Snippet) error {
	args := repo.Called(ctx, snippet)
	return args.Error(0)
}
