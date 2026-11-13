package service

import (
	"context"
	"errors"
	"github.com/AanFLenrir/campuscare-backend/app/model"
	"github.com/AanFLenrir/campuscare-backend/app/repository"
)

// CategoryService provides business logic for Category entity.
type CategoryService interface {
	Create(ctx context.Context, c *model.Category) error
	GetAll(ctx context.Context) ([]*model.Category, error)
	GetByID(ctx context.Context, id string) (*model.Category, error)
	Update(ctx context.Context, c *model.Category) error
	Delete(ctx context.Context, id string) error
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) Create(ctx context.Context, c *model.Category) error {
	if c.Name == "" {
		return errors.New("category name required")
	}
	return s.repo.Create(ctx, c)
}

func (s *categoryService) GetAll(ctx context.Context) ([]*model.Category, error) {
	return s.repo.GetAll(ctx)
}

func (s *categoryService) GetByID(ctx context.Context, id string) (*model.Category, error) {
	if id == "" {
		return nil, errors.New("category id required")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *categoryService) Update(ctx context.Context, c *model.Category) error {
	if c.ID == "" {
		return errors.New("category id required for update")
	}
	return s.repo.Update(ctx, c)
}

func (s *categoryService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("category id required for delete")
	}
	return s.repo.Delete(ctx, id)
}
