package models

// import "gorm.io/gorm"

type User struct {
	UUIDBaseModel
	Name     string `gorm:"type:varchar(100);not null" json:"name"`
	Email    string `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password string `gorm:"type:varchar(100);not null" json:"-"`
}

type SignUpInput struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

type SignInInput struct {
	Email    string `json:"email"  validate:"required"`
	Password string `json:"password"  validate:"required"`
}
