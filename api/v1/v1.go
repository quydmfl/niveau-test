package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Sorting struct {
	SortBy    string `json:"sort_by,omitempty" form:"sort_by" binding:"omitempty" example:"price"`
	SortOrder string `json:"sort_order,omitempty" form:"sort_order" binding:"omitempty,oneof=asc desc" example:"asc"`
}

type Page struct {
	Page int `json:"page" form:"page" binding:"required,gte=1" example:"1"`
	Size int `json:"size" form:"size" binding:"required,gte=10,lte=100" example:"10"`
}

type Pagination struct {
	Page
	TotalRows  int `json:"total_rows" example:"100"`
	TotalPages int `json:"total_pages" example:"10"`
}

func (s *Sorting) SetDefault() {
	if s.SortBy == "" {
		s.SortBy = "added_date"
	}

	if s.SortOrder == "" {
		s.SortOrder = "desc"
	}
}

func HandleSuccess(ctx *gin.Context, data interface{}) {
	if data == nil {
		data = map[string]interface{}{}
	}
	resp := Response{Code: errorCodeMap[ErrSuccess], Message: ErrSuccess.Error(), Data: data}
	if _, ok := errorCodeMap[ErrSuccess]; !ok {
		resp = Response{Code: 0, Message: "", Data: data}
	}
	ctx.JSON(http.StatusOK, resp)
}

func HandleError(ctx *gin.Context, httpCode int, err error, data interface{}) {
	if data == nil {
		data = map[string]string{}
	}
	resp := Response{Code: errorCodeMap[err], Message: err.Error(), Data: data}
	if _, ok := errorCodeMap[err]; !ok {
		resp = Response{Code: 500, Message: "unknown error", Data: data}
	}
	ctx.JSON(httpCode, resp)
}

type Error struct {
	Code    int
	Message string
}

var errorCodeMap = map[error]int{}

func newError(code int, msg string) error {
	err := errors.New(msg)
	errorCodeMap[err] = code
	return err
}
func (e Error) Error() string {
	return e.Message
}
