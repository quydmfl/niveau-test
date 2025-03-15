package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	v1 "github.com/quydmfl/niveau-test/api/v1"
	"github.com/quydmfl/niveau-test/internal/model"
	"github.com/quydmfl/niveau-test/internal/repository"
)

type CategoryService interface {
	GetCategories(ctx context.Context, req *v1.SearchCategoryRequest) (*v1.SearchCategoryResponse, error)
	GetCategory(ctx context.Context, id string) (*v1.GetCategoryDetailData, error)
	CreateCategory(ctx context.Context, req *v1.CreateCategoryRequest) error
}

func NewCategoryService(
	service *Service,
	categoryRepository repository.CategoryRepository,
) CategoryService {
	return &categoryService{
		Service:            service,
		categoryRepository: categoryRepository,
	}
}

type categoryService struct {
	*Service
	categoryRepository repository.CategoryRepository
}

func (s *categoryService) GetCategories(ctx context.Context, req *v1.SearchCategoryRequest) (*v1.SearchCategoryResponse, error) {
	var result []v1.GetCategoryDetailData

	categories, total, err := s.categoryRepository.Search(ctx, req)
	if err != nil {
		return nil, err
	}

	for _, category := range *categories {
		result = append(result, v1.GetCategoryDetailData{
			ID:        category.ID.String(),
			Name:      category.Name,
			Status:    category.Status,
			CreatedAt: category.CreatedAt.Format("2006-01-02"),
			UpdatedAt: category.UpdatedAt.Format("2006-01-02"),
		})
	}

	// Calculate total pages
	totalPages := int(total / req.Size)
	if total%req.Size > 0 {
		totalPages++
	}

	return &v1.SearchCategoryResponse{
		Pagination: v1.Pagination{
			Page:       req.Page,
			TotalRows:  total,
			TotalPages: totalPages,
		},
		Response: v1.Response{
			Data: result,
		},
	}, nil
}

func (s *categoryService) GetCategory(ctx context.Context, id string) (*v1.GetCategoryDetailData, error) {
	categoryUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	category, err := s.categoryRepository.GetCategoryById(ctx, categoryUUID)
	if err != nil {
		return nil, err
	}

	return &v1.GetCategoryDetailData{
		ID:        category.ID.String(),
		Name:      category.Name,
		Status:    category.Status,
		CreatedAt: category.CreatedAt.Format("2006-01-02"),
		UpdatedAt: category.UpdatedAt.Format("2006-01-02"),
	}, nil
}

func (s *categoryService) CreateCategory(ctx context.Context, req *v1.CreateCategoryRequest) error {
	category := &model.Category{
		Name:      req.Name,
		Status:    req.Status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := s.tm.Transaction(ctx, func(ctx context.Context) error {
		if err := s.categoryRepository.Create(ctx, category); err != nil {
			return err
		}
		return nil
	})

	return err
}
