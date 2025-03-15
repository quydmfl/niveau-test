package service

import (
	"context"

	"github.com/google/uuid"
	v1 "github.com/quydmfl/niveau-test/api/v1"
	"github.com/quydmfl/niveau-test/internal/model"
	"github.com/quydmfl/niveau-test/internal/repository"
)

type SupplierService interface {
	GetSuppliers(ctx context.Context, req *v1.SearchSupplierRequest) (*v1.SearchSupplierResponse, error)
	GetSupplier(ctx context.Context, id string) (*v1.GetSupplierDetailData, error)
	CreateSupplier(ctx context.Context, req *v1.CreateSupplierRequest) error
}

func NewSupplierService(
	service *Service,
	supplierRepository repository.SupplierRepository,
) SupplierService {
	return &supplierService{
		Service:            service,
		supplierRepository: supplierRepository,
	}
}

type supplierService struct {
	*Service
	supplierRepository repository.SupplierRepository
}

func (s *supplierService) GetSupplier(ctx context.Context, id string) (*v1.GetSupplierDetailData, error) {
	supplierUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	supplier, err := s.supplierRepository.GetSupplierById(ctx, supplierUUID)
	if err != nil {
		return nil, err
	}

	return &v1.GetSupplierDetailData{
		ID:   supplier.ID.String(),
		Name: supplier.Name,
	}, nil
}

func (s *supplierService) GetSuppliers(ctx context.Context, req *v1.SearchSupplierRequest) (*v1.SearchSupplierResponse, error) {
	var result []v1.GetSupplierDetailData

	suppliers, total, err := s.supplierRepository.Search(ctx, req)
	if err != nil {
		return nil, err
	}

	for _, supplier := range *suppliers {
		result = append(result, v1.GetSupplierDetailData{
			ID:   supplier.ID.String(),
			Name: supplier.Name,
		})
	}

	// Calculate total pages
	totalPages := int(total / req.Size)
	if total%req.Size > 0 {
		totalPages++
	}

	return &v1.SearchSupplierResponse{
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

func (s *supplierService) CreateSupplier(ctx context.Context, req *v1.CreateSupplierRequest) error {
	supplier := &model.Supplier{
		Name: req.Name,
	}

	err := s.tm.Transaction(ctx, func(ctx context.Context) error {
		if err := s.supplierRepository.Create(ctx, supplier); err != nil {
			return err
		}
		return nil
	})

	return err
}
