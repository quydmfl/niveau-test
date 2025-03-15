package repository

import (
	"context"
	"errors"
	"fmt"

	v1 "github.com/quydmfl/niveau-test/api/v1"
	"github.com/quydmfl/niveau-test/internal/model"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(ctx context.Context, product *model.Product) error
	Update(ctx context.Context, product *model.Product) error
	Delete(ctx context.Context, id string) error
	GetProductByPref(ctx context.Context, id string) (*model.Product, error)
	Search(ctx context.Context, req *v1.SearchProductRequest) (*[]model.Product, int, error)
	SumQuantityProducts(ctx context.Context) (int64, error)
	StatsProductsPerCategory(ctx context.Context, totalQuantity int64) ([]v1.ProductCategoryStatsResponse, error)
	StatsProductsPerSupplier(ctx context.Context, totalQuantity int64) ([]v1.ProductSupplierStatsResponse, error)
}

func NewProductRepository(
	repository *Repository,
) ProductRepository {
	return &productRepository{
		Repository: repository,
	}
}

type productRepository struct {
	*Repository
}

func (r *productRepository) Create(ctx context.Context, product *model.Product) error {
	if err := r.DB(ctx).Create(product).Error; err != nil {
		return err
	}
	return nil
}

func (r *productRepository) Update(ctx context.Context, product *model.Product) error {
	if err := r.DB(ctx).Save(product).Error; err != nil {
		return err
	}
	return nil
}

func (r *productRepository) Delete(ctx context.Context, id string) error {
	db := r.DB(ctx)

	// Check if product exists
	var product model.Product
	if err := db.Where("reference = ?", id).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return v1.ErrNotFound
		}
		return err
	}

	// Perform delete
	result := db.Where("reference = ?", id).Delete(&model.Product{})
	if result.Error != nil {
		return result.Error
	}

	// Check if rows were affected
	if result.RowsAffected == 0 {
		return fmt.Errorf("no product deleted, reference %s may not exist", id)
	}

	return nil
}

func (r *productRepository) GetProductByPref(ctx context.Context, id string) (*model.Product, error) {
	var product model.Product
	if err := r.DB(ctx).Preload("Category").Preload("Supplier").Where("reference = ?", id).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, v1.ErrNotFound
		}
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) Search(ctx context.Context, req *v1.SearchProductRequest) (*[]model.Product, int, error) {
	var (
		products  []model.Product
		totalRows int64
	)

	query := r.DB(ctx).Model(&model.Product{})

	// Apply filters
	if req.Reference != "" {
		query = query.Where("reference = ?", req.Reference)
	}
	if req.ProductName != "" {
		query = query.Where("product_name ILIKE ?", "%"+req.ProductName+"%")
	}
	if req.CategoryId != "" {
		query = query.Where("category_id = ?", req.CategoryId)
	}
	if req.SupplierId != "" {
		query = query.Where("supplier_id = ?", req.SupplierId)
	}
	if req.StockLocationId != "" {
		query = query.Where("stock_location_id = ?", req.StockLocationId)
	}
	if req.MinPrice > 0 {
		query = query.Where("price >= ?", req.MinPrice)
	}
	if req.MaxPrice > 0 {
		query = query.Where("price <= ?", req.MaxPrice)
	}
	if req.DateAddedFrom != "" {
		query = query.Where("added_date >= ?", req.DateAddedFrom)
	}
	if req.DateAddedTo != "" {
		query = query.Where("added_date <= ?", req.DateAddedTo)
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
	if err := query.Preload("Category").Preload("Supplier").Find(&products).Error; err != nil {
		return nil, 0, err
	}
	return &products, int(totalRows), nil
}

func (r *productRepository) SumQuantityProducts(ctx context.Context) (int64, error) {
	var count int64
	if err := r.DB(ctx).Model(&model.Product{}).Select("SUM(quantity)").Scan(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *productRepository) StatsProductsPerCategory(ctx context.Context, totalQuantity int64) ([]v1.ProductCategoryStatsResponse, error) {
	var categoryStats []v1.ProductCategoryStatsResponse
	if err := r.DB(ctx).Table("products").
		Select("product_categories.id as category_id, product_categories.name as category_name, ROUND((SUM(products.quantity) * 100.0 / ?)) as percentage", totalQuantity).
		Joins("JOIN product_categories ON product_categories.id = products.category_id").
		Group("product_categories.id, product_categories.name").
		Scan(&categoryStats).Error; err != nil {
		return nil, err
	}

	return categoryStats, nil
}

func (r *productRepository) StatsProductsPerSupplier(ctx context.Context, totalQuantity int64) ([]v1.ProductSupplierStatsResponse, error) {
	var supplierStats []v1.ProductSupplierStatsResponse
	if err := r.DB(ctx).Table("products").
		Select("suppliers.id as supplier_id, suppliers.name as supplier_name, ROUND((SUM(products.quantity) * 100.0 / ?)) as percentage", totalQuantity).
		Joins("JOIN suppliers ON suppliers.id = products.supplier_id").
		Group("suppliers.id, suppliers.name").
		Scan(&supplierStats).Error; err != nil {
		return nil, err
	}

	return supplierStats, nil
}
