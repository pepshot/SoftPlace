package dto

type FurnitureRequest struct {
	Name    string                `json:"name" binding:"required,max=100"`
	Code    string                `json:"code" binding:"required,max=20"`
	Modules []FurnitureModuleItem `json:"modules" binding:"required,min=1"`
}

type FurnitureModuleItem struct {
	ModuleID string `json:"moduleID" binding:"required"`
	Count    int    `json:"count" binding:"required,min=1"`
}

type FurnitureResponse struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Code       string  `json:"code"`
	Price      float64 `json:"price"`
	StockCount int     `json:"stockCount"`
}
