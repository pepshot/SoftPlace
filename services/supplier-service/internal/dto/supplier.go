package dto

type RegisterSupplierRequest struct {
	Login           string `json:"login" binding:"required,min=3,max=50"`
	Password        string `json:"password" binding:"required,min=6,max=50"`
	ConfirmPassword string `json:"confirmPassword" binding:"required"`
	Email           string `json:"email" binding:"required,email,max=100"`
}

type LoginSupplierRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SupplierResponse struct {
	ID    string `json:"id"`
	Login string `json:"login"`
	Email string `json:"email"`
}

type AuthSupplierResponse struct {
	Token    string           `json:"token"`
	Supplier SupplierResponse `json:"supplier"`
}
