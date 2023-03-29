package models

type SignUpRequest struct {
	FirstName   string  `json:"first_name" binding:"required,min=1,max=255"`
	LastName    string  `json:"last_name" binding:"required,min=1,max=255"`
	MiddleName  *string `json:"middle_name" binding:"omitempty,min=1,max=255"`
	CompanyName string  `json:"company_name" binding:"required,min=1,max=255"`
	Email       string  `json:"email" binding:"email,max=255"`
	Phone       string  `json:"phone" binding:"phone"`
	Password    string  `json:"password" binding:"password"`
	LocalISO    string  `json:"local_iso" binding:"required,min=2,max=32"`
}
