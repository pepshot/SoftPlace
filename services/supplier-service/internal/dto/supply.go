package dto

type SupplyRequest struct {
	Code      string               `json:"code" binding:"required,max=20"`
	Date      string               `json:"date" binding:"required"`
	Materials []SupplyMaterialItem `json:"materials" binding:"required,min=1"`
}

type SupplyMaterialItem struct {
	MaterialID string `json:"materialID" binding:"required"`
	Count      int    `json:"count" binding:"required,min=1"`
}

type SupplyResponse struct {
	ID         string  `json:"id"`
	SupplierID string  `json:"supplierID"`
	Code       string  `json:"code"`
	Date       string  `json:"date"`
	Price      float64 `json:"price"`
}
