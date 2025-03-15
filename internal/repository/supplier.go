package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	v1 "github.com/quydmfl/niveau-test/api/v1"
	"github.com/quydmfl/niveau-test/internal/model"
	"gorm.io/gorm"
)

type SupplierRepository interface {
	Create(ctx context.Context, category *model.Supplier) error
	GetSupplierById(ctx context.Context, id uuid.UUID) (*model.Supplier, error)
	Search(ctx context.Context, req *v1.SearchSupplierRequest) (*[]model.Supplier, int, error)
}

func NewSupplierRepository(
	repository *Repository,
) SupplierRepository {
	return &supplierRepository{
		Repository: repository,
	}
}

type supplierRepository struct {
	*Repository
}

func (r *supplierRepository) Create(ctx context.Context, supplier *model.Supplier) error {
	if err := r.DB(ctx).Create(supplier).Error; err != nil {
		return err
	}
	return nil
}

func (r *supplierRepository) GetSupplierById(ctx context.Context, id uuid.UUID) (*model.Supplier, error) {
	var supplier model.Supplier

	if err := r.DB(ctx).Where("id = ?", id).First(&supplier).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, v1.ErrNotFound
		}
		return nil, err
	}

	return &supplier, nil
}

func (r *supplierRepository) Search(ctx context.Context, req *v1.SearchSupplierRequest) (*[]model.Supplier, int, error) {
	var (
		suppliers []model.Supplier
		totalRows int64
	)

	query := r.DB(ctx).Model(&model.Supplier{})

	// Apply filters
	if req.Name != "" {
		query = query.Where("name ILIKE ?", "%"+req.Name+"%")
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

	if err := query.Find(&suppliers).Error; err != nil {
		return nil, 0, err
	}

	return &suppliers, int(totalRows), nil
}
