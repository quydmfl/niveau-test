package handler

import (
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	v1 "github.com/quydmfl/niveau-test/api/v1"
	"github.com/quydmfl/niveau-test/internal/helper"
	"github.com/quydmfl/niveau-test/internal/service"
)

type ProductHandler struct {
	*Handler
	productService service.ProductService
}

func NewProductHandler(
	handler *Handler,
	productService service.ProductService,
) *ProductHandler {
	return &ProductHandler{
		Handler:        handler,
		productService: productService,
	}
}

// GetProducts godoc
// @Summary Search and filter products
// @Description Retrieve a list of products with optional filters, sorting, and pagination
// @Tags Product Modules
// @Accept json
// @Produce json
// @Security Bearer
// @Param reference query string false "Filter by Product Reference" example("PROD-202401-001")
// @Param product_name query string false "Filter by Product Name (3-100 chars)" example("Crunchy Munch")
// @Param category_id query string false "Filter by Category ID (must be > 0)" example(1)
// @Param supplier_id query string false "Filter by Supplier ID (must be > 0)" example(1)
// @Param stock_location_id query string false "Filter by Stock Location ID (must be > 0)" example(1)
// @Param min_price query number false "Filter by Minimum Price (must be > 0)" example(100000)
// @Param max_price query number false "Filter by Maximum Price (must be > 0)" example(500000)
// @Param date_added_from query string false "Filter by start date (YYYY-MM-DD)" format(date) example("2024-01-28")
// @Param date_added_to query string false "Filter by end date (YYYY-MM-DD)" format(date) example("2024-01-28")
// @Param status query string false "Filter by Product Status" Enums(Available, Out of Stock, On Order) example("Available")
// @Param sort_by query string false "Sort by field" Enums(price, name, added_date) default(added_date) example("price")
// @Param sort_order query string false "Sort order" Enums(asc, desc) default(desc) example("asc")
// @Param page query int true "Page number (must be >= 1)" default(1) example(1)
// @Param size query int true "Items per page (between 10-100)" default(20) example(10)
// @Success 200 {object} v1.SearchProductResponse
// @Router /products [get]
func (h *ProductHandler) GetProducts(ctx *gin.Context) {
	var req v1.SearchProductRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	// Set default value for sorting
	req.Sorting.SetDefault()

	products, err := h.productService.SearchProduct(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}

	v1.HandleSuccess(ctx, products)
}

// GetProductDetail godoc
// @Summary Get product detail information
// @Description Retrieve detailed product information by ID
// @Tags Product Modules
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Product Reference ID"
// @Success 200 {object} v1.GetProductDetailData
// @Router /products/{id} [get]
func (h *ProductHandler) GetProductDetail(ctx *gin.Context) {
	userId := GetUserIdFromCtx(ctx)
	if userId == "" {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		return
	}

	productPref := ctx.Param("id")
	if productPref == "" {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	product, err := h.productService.GetProduct(ctx, productPref)

	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	v1.HandleSuccess(ctx, product)
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product with the given details
// @Tags Product Modules
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.CreateProductRequest true "Product creation request"
// @Success 201 {object} v1.Response "Product created successfully"
// @Router /products [post]
func (h *ProductHandler) CreateProduct(ctx *gin.Context) {
	userId := GetUserIdFromCtx(ctx)
	if userId == "" {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		return
	}

	req := new(v1.CreateProductRequest)
	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.productService.CreateProduct(ctx, userId, req); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}

// UpdateProduct godoc
// @Summary Update a product by ID
// @Description Update product details using the product reference (ID) provided in the URL
// @Tags Product Modules
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Product Reference ID"
// @Param body body v1.UpdateProductRequest true "Update Product Request Body"
// @Success 200 {object} map[string]interface{} "Success message"
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(ctx *gin.Context) {
	userId := GetUserIdFromCtx(ctx)
	productPref := ctx.Param("id")
	if productPref == "" {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	var req v1.UpdateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.productService.UpdateProduct(ctx, userId, productPref, &req); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// DeleteProduct godoc
// @Summary Delete a product by ID
// @Description Permanently deletes a product using the provided product reference (ID) in the URL
// @Tags Product Modules
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Product Reference ID"
// @Success 200 {object} map[string]interface{} "Product deleted successfully"
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(ctx *gin.Context) {
	productPref := ctx.Param("id")
	if productPref == "" {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.productService.DeleteProduct(ctx, productPref); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}

// GetProductsPerCategory godoc
// @Summary Get percentage of products per category
// @Description Retrieve statistics on the distribution of products per category based on quantity
// @Tags Statistics
// @Produce json
// @Security Bearer
// @Success 200 {array} v1.ProductCategoryStatsResponse
// @Router /statistics/products-per-category [get]
func (h *ProductHandler) GetProductsPerCategory(ctx *gin.Context) {
	var categoryStats []v1.ProductCategoryStatsResponse

	categoryStats, err := h.productService.StatsProductsPerCategory(ctx)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}

	v1.HandleSuccess(ctx, categoryStats)
}

// GetProductsPerSupplier godoc
// @Summary Get percentage of products per supplier
// @Description Retrieve statistics on the distribution of products per supplier based on quantity
// @Tags Statistics
// @Produce json
// @Security Bearer
// @Success 200 {array} v1.ProductSupplierStatsResponse
// @Router /statistics/products-per-supplier [get]
func (h *ProductHandler) GetProductsPerSupplier(ctx *gin.Context) {
	var supplierStats []v1.ProductSupplierStatsResponse

	supplierStats, err := h.productService.StatsProductsPerSupplier(ctx)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}

	v1.HandleSuccess(ctx, supplierStats)
}

// ExportProducts godoc
// @Summary Export a single product
// @Description Exports product details in the specified format (currently supports PDF).
// @Produce application/json
// @Tags Product Modules
// @Accept json
// @Produce json
// @Security Bearer
// @Param format path string true "Export format (e.g., 'pdf')"
// @Param id path string true "Product ID (UUID)"
// @Success 200 {object} v1.Response "File generated successfully"
// @Router /products/export/{format}/{id} [get]
func (h *ProductHandler) ExportProducts(ctx *gin.Context) {
	var err error
	format := ctx.Param("format")
	productId := ctx.Param("id")

	switch format {
	case "pdf":
		err = h.productService.ExportProductsToPDF(ctx, productId)
	default:
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}

// GetDistanceBetweenIPAndCity godoc
// @Summary Calculate the distance between an IP address and a city
// @Description Computes the distance (in km) between the user's location (determined via IP) and a specified city.
// @Tags Product Modules
// @Accept json
// @Produce json
// @Security Bearer
// @Param city path string true "City name to compare distance with"
// @Success 200 {object} map[string]interface{} "Returns the calculated distance in kilometers"
// @Router /products/distance/ip/{city} [get]
func (h *ProductHandler) GetDistanceBetweenIPAndCity(ctx *gin.Context) {
	// Get client's IP automatically
	ip := helper.GetRealIP(ctx)
	city := ctx.Param("city")

	from, to, distance, err := h.productService.CalculateDistance(ctx, ip, city)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}

	v1.HandleSuccess(ctx, map[string]interface{}{"from": from, "to": to, "distance": math.Round(distance), "unit": "km"})
}
