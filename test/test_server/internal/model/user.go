package model

type User struct {
	Name     string   `json:"name" validate:"required"`
	Username string   `json:"username" validate:"required,username,min=5,max=10"`
	Email    string   `json:"email" validate:"required,email"`
	Password string   `json:"password" validate:"required,min=8"`
	Roles    []string `gorm:"many2many:user_role;"`
}
