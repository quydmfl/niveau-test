package v1

type CreateSupplierRequest struct {
	Name string `json:"name" form:"name" binding:"omitempty,min=3,max=100" example:"Crunchy Munch"`
}

type SearchSupplierRequest struct {
	Page
	Sorting

	Name string `json:"name" form:"name" binding:"omitempty,min=3,max=100" example:"Crunchy Munch"`
}

type SearchSupplierResponse struct {
	Response
	Pagination
}

type GetSupplierDetailData struct {
	ID   string `json:"id" example:"67389c69-2e78-413a-9d77-6b749520b127"`
	Name string `json:"name" example:"Food"`
}

type GetSupplierDetailResponse struct {
	Response
	Data GetSupplierDetailData
}
