package recipe

import (
	"context"
	"foodplanner/internal/db"
)

type Service struct {
	db   db.DBTX
	Repo *Repo
}

func NewService(db db.DBTX, repo *Repo) *Service {
	return &Service{
		db:   db,
		Repo: repo,
	}
}

func (s *Service) CreateRecipe(ctx context.Context, name string) (*Recipe, error) {
	entity, err := NewRecipe(name)
	if err != nil {
		return nil, err
	}
	return s.Repo.CreateRecipe(ctx, s.db, entity)
}

func (s *Service) GetAllRecipes(ctx context.Context) ([]*Recipe, error) {
	return s.Repo.GetAllRecipes(ctx, s.db)
}

func (s *Service) GetRecipeByID(ctx context.Context, id string) (*Recipe, error) {
	return s.Repo.GetRecipeByID(ctx, s.db, id)
}
