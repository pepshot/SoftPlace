package dto

type MaterialRequest struct {
	Name       string  `json:"name" binding:"required,max=100"`
	Code       string  `json:"code" binding:"required,max=20"`
	Price      float64 `json:"price" binding:"required,min=0"`
	StockCount int     `json:"stockCount" binding:"omitempty,min=0"`
}

type MaterialResponse struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Code       string  `json:"code"`
	Price      float64 `json:"price"`
	StockCount int     `json:"stockCount"`
}
