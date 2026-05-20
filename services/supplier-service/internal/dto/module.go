package dto

type ModuleRequest struct {
	Name      string               `json:"name" binding:"required,max=100"`
	Code      string               `json:"code" binding:"required,max=20"`
	Materials []ModuleMaterialItem `json:"materials" binding:"required,min=1"`
}

type ModuleResponse struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Code       string  `json:"code"`
	Price      float64 `json:"price"`
	StockCount int     `json:"stockCount"`
}

type ModuleItem struct {
	ModuleID string `json:"moduleID" binding:"required"`
	Count    int    `json:"count" binding:"required,min=1"`
}

type ModuleMaterialItem struct {
	MaterialID string `json:"materialID" binding:"required"`
	Count      int    `json:"count" binding:"required,min=1"`
}
