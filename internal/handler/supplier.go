package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	v1 "github.com/quydmfl/niveau-test/api/v1"
	"github.com/quydmfl/niveau-test/internal/service"
)

type SupplierHandler struct {
	*Handler
	supplierService service.SupplierService
}

func NewSupplierHandler(
	handler *Handler,
	supplierService service.SupplierService,
) *SupplierHandler {
	return &SupplierHandler{
		Handler:         handler,
		supplierService: supplierService,
	}
}

// GetSuppliers godoc
//
// @Summary Get list of suppliers
// @Description Retrieve a paginated list of suppliers with optional search and sorting criteria.
// @Tags Supplier Modules
// @Accept json
// @Produce json
// @Security Bearer
// @Param name query string false "Filter by supplier name"
// @Param sort_by query string false "Sort by field" Enums(name, id) default(name) example("name")
// @Param sort_order query string false "Sort order" Enums(asc, desc) default(desc) example("asc")
// @Param page query int true "Page number (must be >= 1)" default(1) example(1)
// @Param size query int true "Items per page (between 10-100)" default(20)  example(10)
// @Success 200 {array} v1.SearchSupplierResponse "Successful response with a list of suppliers"
// @Router /supplier [get]
func (h *SupplierHandler) GetSuppliers(ctx *gin.Context) {
	var req v1.SearchSupplierRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	// Set default value for sorting
	req.Sorting.SetDefault()

	suppliers, err := h.supplierService.GetSuppliers(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	v1.HandleSuccess(ctx, suppliers)
}

// GetSupplierDetail godoc.
//
// @Summary Get supplier details
// @Description Retrieve detailed information about a supplier by its ID.
// @Tags Supplier Modules
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Supplier ID"
// @Success 200 {object} v1.GetSupplierDetailResponse "Successful response"
// @Router /supplier/{id} [get]
func (h *SupplierHandler) GetSupplierDetail(ctx *gin.Context) {
	supplierID := ctx.Param("id")

	if supplierID == "" {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	supplier, err := h.supplierService.GetSupplier(ctx, supplierID)
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	v1.HandleSuccess(ctx, supplier)
}

// CreateSupplier godoc
// @Summary Create a new supplier
// @Description Create a new supplier with the given details
// @Tags Supplier Modules
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.CreateSupplierRequest true "Supplier creation request"
// @Success 201 {object} v1.Response "Supplier created successfully"
// @Router /supplier [post]
func (h *SupplierHandler) CreateSupplier(ctx *gin.Context) {
	req := new(v1.CreateSupplierRequest)
	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.supplierService.CreateSupplier(ctx, req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}
