package v1

type SearchCategoryRequest struct {
	Page
	Sorting

	Name   string `json:"name" form:"name" binding:"omitempty,min=3,max=100" example:"Crunchy Munch"`
	Status string `json:"status" form:"status" binding:"omitempty,oneof='active' 'deactive'" example:"active"`
}

type SearchCategoryResponse struct {
	Response
	Pagination
}

type CreateCategoryRequest struct {
	Name   string `json:"name" form:"name" binding:"omitempty,min=3,max=100" example:"Crunchy Munch"`
	Status string `json:"status" form:"status" binding:"omitempty,oneof='active' 'deactive'" example:"active"`
}

type GetCategoryDetailData struct {
	ID        string `json:"id" example:"67389c69-2e78-413a-9d77-6b749520b127"`
	Name      string `json:"name" example:"Food"`
	Status    string `json:"status" example:"active"`
	CreatedAt string `json:"created_at" example:"2023-01-01"`
	UpdatedAt string `json:"updated_at" example:"2023-01-01"`
}

type GetCategoryDetailResponse struct {
	Response
	Data GetCategoryDetailData
}
