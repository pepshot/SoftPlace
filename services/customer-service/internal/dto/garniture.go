package dto

type GarnitureRequest struct {
	Name      string                   `json:"name" binding:"required,max=100"`
	Code      string                   `json:"code" binding:"required,max=20"`
	Furniture []GarnitureFurnitureItem `json:"furniture" binding:"required,min=1"`
}

type GarnitureFurnitureItem struct {
	FurnitureID string `json:"furnitureId" binding:"required"`
	Count       int    `json:"count" binding:"required,min=1"`
}

type GarnitureResponse struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Code       string  `json:"code"`
	Price      float64 `json:"price"`
	StockCount int     `json:"stockCount"`
}
