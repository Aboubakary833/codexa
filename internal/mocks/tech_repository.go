package mocks

import (
	"context"

	"github.com/aboubakary833/codexa/internal/domain"
	"github.com/stretchr/testify/mock"
)

type techRepository struct {
	mock.Mock
}

func NewTechRepository() *techRepository {
	return new(techRepository)
}

func (repo *techRepository) Store(ctx context.Context, tech domain.Tech) error {
	args := repo.Called(ctx, tech)
	return args.Error(0)
}

func (repo *techRepository) FindAll(ctx context.Context) ([]domain.Tech, error) {
	args := repo.Called(ctx)

	if args.Get(0) == nil {
		return []domain.Tech{}, args.Error(1)
	}

	return args.Get(0).([]domain.Tech), nil
}

func (repo *techRepository) Retrieve(ctx context.Context, needle string) (domain.Tech, error) {
	args := repo.Called(ctx, needle)

	if args.Get(0) == nil {
		return domain.Tech{}, args.Error(1)
	}

	return args.Get(0).(domain.Tech), nil
}

func (repo *techRepository) GetAliases(ctx context.Context, techID string) ([]domain.TechAlias, error) {
	args := repo.Called(ctx, techID)

	if args.Get(0) == nil {
		return []domain.TechAlias{}, args.Error(1)
	}

	return args.Get(0).([]domain.TechAlias), nil
}
