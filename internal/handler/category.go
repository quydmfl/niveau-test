package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	v1 "github.com/quydmfl/niveau-test/api/v1"
	"github.com/quydmfl/niveau-test/internal/service"
)

type CategoryHandler struct {
	*Handler
	categoryService service.CategoryService
}

func NewCategoryHandler(
	handler *Handler,
	categoryService service.CategoryService,
) *CategoryHandler {
	return &CategoryHandler{
		Handler:         handler,
		categoryService: categoryService,
	}
}

// GetCategories godoc
//
// @Summary Get list of categories
// @Description Retrieve a paginated list of categories with optional search and sorting criteria.
// @Tags Categories Modules
// @Accept json
// @Produce json
// @Security Bearer
// @Param name query string false "Filter by category name"
// @Param status query string false "Filter by Status" Enums(active, deactive) example("active")
// @Param sort_by query string false "Sort by field" Enums(name, created_at) default(created_at) example("name")
// @Param sort_order query string false "Sort order" Enums(asc, desc) default(desc) example("asc")
// @Param page query int true "Page number (must be >= 1)" default(1) example(1)
// @Param size query int true "Items per page (between 10-100)" default(20)  example(10)
// @Success 200 {array} v1.SearchCategoryResponse "Successful response with a list of categories"
// @Router /categories [get]
func (h *CategoryHandler) GetCategories(ctx *gin.Context) {
	var req v1.SearchCategoryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	// Set default value for sorting
	req.Sorting.SetDefault()

	categories, err := h.categoryService.GetCategories(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	v1.HandleSuccess(ctx, categories)
}

// GetCategoryDetail godoc.
//
// @Summary Get category details
// @Description Retrieve detailed information about a category by its ID.
// @Tags Categories Modules
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Category ID"
// @Success 200 {object} v1.GetCategoryDetailResponse "Successful response"
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetCategoryDetail(ctx *gin.Context) {
	categoryID := ctx.Param("id")

	if categoryID == "" {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	category, err := h.categoryService.GetCategory(ctx, categoryID)
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	v1.HandleSuccess(ctx, category)
}

// CreateCategory godoc
// @Summary Create a new category
// @Description Create a new category with the given details
// @Tags Categories Modules
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.CreateCategoryRequest true "Category creation request"
// @Success 201 {object} v1.Response "Category created successfully"
// @Router /categories [post]
func (h *CategoryHandler) CreateCategory(ctx *gin.Context) {
	req := new(v1.CreateCategoryRequest)
	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.categoryService.CreateCategory(ctx, req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}
