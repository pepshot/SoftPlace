package dto

type ShipmentRequest struct {
	Code       string                  `json:"code" binding:"required,max=20"`
	Date       string                  `json:"date" binding:"required"` // формат: YYYY-MM-DD
	Garnitures []ShipmentGarnitureItem `json:"garnitures" binding:"required,min=1"`
}

type ShipmentGarnitureItem struct {
	GarnitureID string `json:"garnitureId" binding:"required"`
	Count       int    `json:"count" binding:"required,min=1"`
}

type ShipmentResponse struct {
	ID    string  `json:"id"`
	Code  string  `json:"code"`
	Date  string  `json:"date"`
	Price float64 `json:"price"`
}
