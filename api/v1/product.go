package v1

type SearchProductRequest struct {
	Page
	Sorting

	Reference       string  `json:"reference" form:"reference" binding:"omitempty,product_ref" example:"PROD-202401-001"`
	ProductName     string  `json:"product_name" form:"product_name" binding:"omitempty,min=3,max=100" example:"Crunchy Munch"`
	CategoryId      string  `json:"category_id" form:"category_id" binding:"omitempty,gt=0" example:"1"`
	SupplierId      string  `json:"supplier_id" form:"supplier_id" binding:"omitempty,gt=0" example:"1"`
	StockLocationId string  `json:"stock_location_id" form:"stock_location_id" binding:"omitempty,gt=0" example:"1"`
	MinPrice        float64 `json:"min_price" form:"min_price" binding:"omitempty,gt=0" example:"100000"`
	MaxPrice        float64 `json:"max_price" form:"max_price" binding:"omitempty,gt=0" example:"500000"`
	DateAddedFrom   string  `json:"date_added_from" form:"date_added_from" binding:"omitempty,datetime=2006-01-02" example:"2024-01-28"`
	DateAddedTo     string  `json:"date_added_to" form:"date_added_to" binding:"omitempty,datetime=2006-01-02" example:"2024-01-28"`
	Status          string  `json:"status" form:"status" binding:"omitempty,oneof='Available' 'Out of Stock' 'On Order'" example:"Available"`
}

type SearchProductResponse struct {
	Response
	Pagination
}

type CreateProductRequest struct {
	Reference     string  `json:"reference" binding:"omitempty,product_ref" example:"PROD-202401-001"`
	ProductName   string  `json:"product_name" binding:"required,min=3,max=100" example:"Crunchy Munch"`
	CategoryId    string  `json:"category_id" binding:"required,uuid4" example:"67389c69-2e78-413a-9d77-6b749520b127"`
	Price         float64 `json:"price" binding:"required,gt=0" example:"150000"`
	Status        string  `json:"status" binding:"required,oneof='Available' 'Out of Stock' 'On Order'" example:"Available"`
	StockLocation string  `json:"stock_location" binding:"required" example:"Brookstone"`
	SupplierId    string  `json:"supplier_id" binding:"required,uuid4" example:"66ae8b26-a0cd-40d5-8b2d-5127d7b9d817"`
	Quantity      int     `json:"quantity" binding:"required,gt=0" example:"10"`
}

type UpdateProductRequest struct {
	ProductName   string  `json:"product_name" binding:"required,min=3,max=100" example:"Crunchy Munch"`
	CategoryId    string  `json:"category_id" binding:"required,uuid4" example:"67389c69-2e78-413a-9d77-6b749520b127"`
	Price         float64 `json:"price" binding:"required,gt=0" example:"150000"`
	Status        string  `json:"status" binding:"required,oneof='Available' 'Out of Stock' 'On Order'" example:"Available"`
	StockLocation string  `json:"stock_location" binding:"required" example:"Brookstone"`
	DateAdded     string  `json:"added_date" binding:"required,datetime=2006-01-02" example:"2024-01-28"`
	SupplierId    string  `json:"supplier_id" binding:"required,uuid4" example:"66ae8b26-a0cd-40d5-8b2d-5127d7b9d817"`
	Quantity      int     `json:"quantity" binding:"required,gt=0" example:"10"`
}

type GetProductDetailData struct {
	Reference     string  `json:"reference" example:"PROD-202401-001"`
	ProductName   string  `json:"product_name" example:"Crunchy Munch"`
	Category      string  `json:"category" example:"Food"`
	Price         float64 `json:"price" example:"150000"`
	Status        string  `json:"status" example:"Available"`
	StockLocation string  `json:"stock_location" example:"Cedar Valley"`
	DateAdded     string  `json:"added_date" example:"2024-01-28"`
	Quantity      int     `json:"quantity" example:"10"`
	Supplier      string  `json:"supplier" example:"KFC"`
}

type ProductCategoryStatsResponse struct {
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Percentage   float64 `json:"percentage"`
}

type ProductSupplierStatsResponse struct {
	SupplierID   string  `json:"supplier_id"`
	SupplierName string  `json:"supplier_name"`
	Percentage   float64 `json:"percentage"`
}
