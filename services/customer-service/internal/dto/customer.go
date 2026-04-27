package dto

type RegisterCustomerRequest struct {
	Login           string `json:"login" binding:"required,min=3,max=50"`
	Password        string `json:"password" binding:"required,min=6,max=50"`
	ConfirmPassword string `json:"confirmPassword" binding:"required"`
	Email           string `json:"email" binding:"required,email,max=100"`
}

type LoginCustomerRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CustomerResponse struct {
	ID    string `json:"id"`
	Login string `json:"login"`
	Email string `json:"email"`
}

type AuthCustomerResponse struct {
	Token    string           `json:"token"`
	Customer CustomerResponse `json:"customer"`
}
