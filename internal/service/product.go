package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	v1 "github.com/quydmfl/niveau-test/api/v1"
	"github.com/quydmfl/niveau-test/internal/helper"
	"github.com/quydmfl/niveau-test/internal/model"
	"github.com/quydmfl/niveau-test/internal/repository"

	"github.com/jung-kurt/gofpdf"
)

type ProductService interface {
	SearchProduct(ctx context.Context, req *v1.SearchProductRequest) (*v1.SearchProductResponse, error)
	GetProduct(ctx context.Context, productRef string) (*v1.GetProductDetailData, error)
	CreateProduct(ctx context.Context, userId string, req *v1.CreateProductRequest) error
	UpdateProduct(ctx context.Context, userId string, productRef string, req *v1.UpdateProductRequest) error
	DeleteProduct(ctx context.Context, productRef string) error
	ExportProductsToPDF(ctx context.Context, productPref string) error
	CalculateDistance(ctx context.Context, ip string, city string) (from string, to string, distance float64, err error)
	StatsProductsPerCategory(ctx context.Context) ([]v1.ProductCategoryStatsResponse, error)
	StatsProductsPerSupplier(ctx context.Context) ([]v1.ProductSupplierStatsResponse, error)
}

func NewProductService(
	service *Service,
	productRepository repository.ProductRepository,
	categoryRepository repository.CategoryRepository,
	documentRepository repository.DocumentsRepository,
	supplierRepository repository.SupplierRepository,
) ProductService {
	return &productService{
		Service:            service,
		productRepository:  productRepository,
		categoryRepository: categoryRepository,
		documentRepository: documentRepository,
		supplierRepository: supplierRepository,
	}
}

type productService struct {
	*Service
	productRepository  repository.ProductRepository
	documentRepository repository.DocumentsRepository
	categoryRepository repository.CategoryRepository
	supplierRepository repository.SupplierRepository
}

func (s *productService) SearchProduct(ctx context.Context, req *v1.SearchProductRequest) (*v1.SearchProductResponse, error) {
	var result []v1.GetProductDetailData
	products, total, err := s.productRepository.Search(ctx, req)
	if err != nil {
		return nil, err
	}

	for _, product := range *products {
		result = append(result, v1.GetProductDetailData{
			Reference:     product.Reference,
			ProductName:   product.Name,
			Category:      product.Category.Name,
			Price:         product.Price,
			Status:        product.Status,
			StockLocation: product.StockCity,
			DateAdded:     product.DateAdded.Format("2006-01-02"),
			Quantity:      product.Quantity,
			Supplier:      product.Supplier.Name,
		})
	}

	// Calculate total pages
	totalPages := int(total / req.Size)
	if total%req.Size > 0 {
		totalPages++
	}

	return &v1.SearchProductResponse{
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

func (s *productService) GetProduct(ctx context.Context, productRef string) (*v1.GetProductDetailData, error) {
	product, err := s.productRepository.GetProductByPref(ctx, productRef)
	if err != nil {
		return nil, err
	}

	return &v1.GetProductDetailData{
		Reference:     product.Reference,
		ProductName:   product.Name,
		Category:      product.Category.Name,
		Price:         product.Price,
		Status:        product.Status,
		StockLocation: product.StockCity,
		DateAdded:     product.DateAdded.Format("2006-01-02"),
		Quantity:      product.Quantity,
		Supplier:      product.Supplier.Name,
	}, nil
}

func (s *productService) CreateProduct(ctx context.Context, userId string, req *v1.CreateProductRequest) error {
	categoryUUID, _ := uuid.Parse(req.CategoryId)
	if _, err := s.categoryRepository.GetCategoryById(ctx, categoryUUID); err != nil {
		return fmt.Errorf("category %s not found", req.CategoryId)
	}

	supplierUUID, _ := uuid.Parse(req.SupplierId)
	if _, err := s.supplierRepository.GetSupplierById(ctx, supplierUUID); err != nil {
		return fmt.Errorf("supplier %s not found", req.SupplierId)
	}

	product := &model.Product{
		Reference:  req.Reference,
		Name:       req.ProductName,
		CategoryID: categoryUUID,
		SupplierID: supplierUUID,
		StockCity:  req.StockLocation,
		Price:      req.Price,
		Status:     req.Status,
		DateAdded:  time.Now(),
		Quantity:   req.Quantity,
	}

	err := s.tm.Transaction(ctx, func(ctx context.Context) error {
		if err := s.productRepository.Create(ctx, product); err != nil {
			return err
		}
		return nil
	})

	return err
}

func (s *productService) UpdateProduct(ctx context.Context, userId string, productRef string, req *v1.UpdateProductRequest) error {
	// Before validator from request, we're skip check error in here
	dateAdded, err := time.Parse("2006-01-02", req.DateAdded)
	if err != nil {
		return err
	}

	categoryUUID, _ := uuid.Parse(req.CategoryId)
	if _, err := s.categoryRepository.GetCategoryById(ctx, categoryUUID); err != nil {
		return fmt.Errorf("category %s not found", req.CategoryId)
	}

	supplierUUID, _ := uuid.Parse(req.SupplierId)
	if _, err := s.supplierRepository.GetSupplierById(ctx, supplierUUID); err != nil {
		return fmt.Errorf("supplier %s not found", req.SupplierId)
	}

	product, err := s.productRepository.GetProductByPref(ctx, productRef)
	if err != nil {
		return err
	}

	product.Name = req.ProductName
	product.CategoryID = categoryUUID
	product.Price = req.Price
	product.Status = req.Status
	product.DateAdded = dateAdded
	product.SupplierID = supplierUUID
	product.Quantity = req.Quantity

	if err = s.productRepository.Update(ctx, product); err != nil {
		return err
	}

	return nil
}

func (s *productService) DeleteProduct(ctx context.Context, productRef string) error {
	return s.productRepository.Delete(ctx, productRef)
}

func (s *productService) ExportProductsToPDF(ctx context.Context, productPref string) error {
	product, err := s.productRepository.GetProductByPref(ctx, productPref)
	if err != nil {
		return err
	}

	fileName, filePath, err := exportProductPdf(ctx, product)
	if err != nil {
		return err
	}

	// Store file information in the database //
	document := model.Documents{
		Filename:   fileName,
		Path:       filePath,
		ProductID:  &product.ID,
		UploadedAt: time.Now(),
	}

	if err := s.documentRepository.Create(ctx, &document); err != nil {
		return err
	}

	return nil
}

func exportProductPdf(ctx context.Context, product *model.Product) (string, string, error) {
	storagePath := "./storage/pdf"
	err := os.MkdirAll(storagePath, os.ModePerm)
	if err != nil {
		return "", "", err
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	// Heading //
	pdf.Cell(40, 10, "Product Information")
	pdf.Ln(12)

	// Product Details //
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("UUID: %s", product.ID.String()))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Reference: %s", product.Reference))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Name: %s", product.Name))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Status: %s", product.Status))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Category Name: %s", product.Category.Name))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Price: $%.2f", product.Price))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Quantity: %d", product.Quantity))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Supplier Name: %s", product.Supplier.Name))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Date Added: %s", product.DateAdded.Format("2006-01-02")))
	pdf.Ln(8)

	// Unique file name //
	fileName := fmt.Sprintf("product_%s_%s.pdf", product.Reference, time.Now().Format("20060102_150405"))
	filePath := filepath.Join(storagePath, fileName)

	// Save the PDF file //
	if err := pdf.OutputFileAndClose(filePath); err != nil {
		return "", "", err
	}

	return fileName, filePath, nil
}

func (s *productService) CalculateDistance(ctx context.Context, ip, city string) (string, string, float64, error) {
	// Get location of the IP
	geoLocation, err := helper.GetGeoLocationByIP(ip)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to get geolocation: %w", err)
	}

	// Get location of the city
	// Currently, I'm using third party API. You can use database to avoid external API
	cityLocation, err := helper.GetGeoLocationByCity(city)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to get geolocation: %w", err)
	}

	// calculate distance //
	return geoLocation.Name, cityLocation.Name, helper.CalculateDistance(geoLocation.Latitude, geoLocation.Longitude, cityLocation.Latitude, cityLocation.Longitude), nil
}

func (s *productService) StatsProductsPerCategory(ctx context.Context) ([]v1.ProductCategoryStatsResponse, error) {
	var (
		totalQuantity int64
		categoryStats []v1.ProductCategoryStatsResponse
		err           error
	)

	totalQuantity, err = s.productRepository.SumQuantityProducts(ctx)
	if err != nil {
		return nil, err
	}

	categoryStats, err = s.productRepository.StatsProductsPerCategory(ctx, totalQuantity)
	if err != nil {
		return nil, err
	}

	return categoryStats, nil
}

func (s *productService) StatsProductsPerSupplier(ctx context.Context) ([]v1.ProductSupplierStatsResponse, error) {
	var (
		totalQuantity int64
		supplierStats []v1.ProductSupplierStatsResponse
		err           error
	)

	totalQuantity, err = s.productRepository.SumQuantityProducts(ctx)
	if err != nil {
		return nil, err
	}

	supplierStats, err = s.productRepository.StatsProductsPerSupplier(ctx, totalQuantity)
	if err != nil {
		return nil, err
	}

	return supplierStats, nil
}
