package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	v1 "github.com/quydmfl/niveau-test/api/v1"
	"github.com/quydmfl/niveau-test/internal/model"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *model.Category) error
	GetCategoryById(ctx context.Context, id uuid.UUID) (*model.Category, error)
	Search(ctx context.Context, req *v1.SearchCategoryRequest) (*[]model.Category, int, error)
}

func NewCategoryRepository(
	repository *Repository,
) CategoryRepository {
	return &categoryRepository{
		Repository: repository,
	}
}

type categoryRepository struct {
	*Repository
}

func (r *categoryRepository) Create(ctx context.Context, category *model.Category) error {
	if err := r.DB(ctx).Create(category).Error; err != nil {
		return err
	}
	return nil
}

func (r *categoryRepository) GetCategoryById(ctx context.Context, id uuid.UUID) (*model.Category, error) {
	var category model.Category

	if err := r.DB(ctx).Where("id = ?", id).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, v1.ErrNotFound
		}
		return nil, err
	}

	return &category, nil
}

func (r *categoryRepository) Search(ctx context.Context, req *v1.SearchCategoryRequest) (*[]model.Category, int, error) {
	var (
		categories []model.Category
		totalRows  int64
	)

	query := r.DB(ctx).Model(&model.Category{})

	// Apply filters
	if req.Name != "" {
		query = query.Where("name ILIKE ?", "%"+req.Name+"%")
	}

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	// Apply sorting
	query = query.Order(req.SortBy + " " + req.SortOrder)

	// Get total count before pagination
	if err := query.Count(&totalRows).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	if req.Page.Page > 0 && req.Size > 0 {
		query = query.Offset((req.Page.Page - 1) * req.Size).Limit(req.Size)
	}

	// Fetch results
	if err := query.Find(&categories).Error; err != nil {
		return nil, 0, err
	}

	return &categories, int(totalRows), nil
}
