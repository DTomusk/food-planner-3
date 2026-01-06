package recipe

import (
	"context"
	"food-planner/internal/db"
)

type Service struct {
	Repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{
		Repo: repo,
	}
}

func (s *Service) CreateRecipe(name string, ctx context.Context, db db.DBTX) (*Recipe, error) {
	entity, err := NewRecipe(name)
	if err != nil {
		return nil, err
	}
	return s.Repo.CreateRecipe(entity, ctx, db)
}
