package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	v1 "github.com/quydmfl/niveau-test/api/v1"
	"github.com/quydmfl/niveau-test/internal/service"
	"go.uber.org/zap"
)

type UserHandler struct {
	*Handler
	userService service.UserService
}

func NewUserHandler(handler *Handler, userService service.UserService) *UserHandler {
	return &UserHandler{
		Handler:     handler,
		userService: userService,
	}
}

// Register godoc
// @Summary User Registration
// @Schemes
// @Description Currently only supports email login
// @Tags Auth Modules
// @Accept json
// @Produce json
// @Param request body v1.RegisterRequest true "params"
// @Success 200 {object} v1.Response
// @Router /auth/register [post]
func (h *UserHandler) Register(ctx *gin.Context) {
	req := new(v1.RegisterRequest)
	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.userService.Register(ctx, req); err != nil {
		h.logger.WithContext(ctx).Error("userService.Register error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}

// Login godoc
// @Summary Account Login
// @Schemes
// @Description
// @Tags Auth Modules
// @Accept json
// @Produce json
// @Param request body v1.LoginRequest true "params"
// @Success 200 {object} v1.LoginResponse
// @Router /auth/login [post]
func (h *UserHandler) Login(ctx *gin.Context) {
	var req v1.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	token, err := h.userService.Login(ctx, &req)
	if err != nil {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		return
	}
	v1.HandleSuccess(ctx, v1.LoginResponseData{
		AccessToken: token,
	})
}

// GetProfile godoc
// @Summary Get user information
// @Schemes
// @Description
// @Tags User Modules
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} v1.GetProfileResponse
// @Router /user [get]
func (h *UserHandler) GetProfile(ctx *gin.Context) {
	userId := GetUserIdFromCtx(ctx)
	if userId == "" {
		v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
		return
	}

	user, err := h.userService.GetProfile(ctx, userId)
	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	v1.HandleSuccess(ctx, user)
}

// UpdateProfile godoc
// @Summary Modify user information
// @Schemes
// @Description
// @Tags User Modules
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body v1.UpdateProfileRequest true "params"
// @Success 200 {object} v1.Response
// @Router /user [put]
func (h *UserHandler) UpdateProfile(ctx *gin.Context) {
	userId := GetUserIdFromCtx(ctx)

	var req v1.UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.userService.UpdateProfile(ctx, userId, &req); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}
