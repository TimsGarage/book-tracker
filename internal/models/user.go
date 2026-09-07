package models

type User struct {
	Model
	Username string `gorm:"not null,unique" json:"username"`
	Password string `gorm:"not null" json:"-"`

	Books []Book `gorm:"foreignKey:UserId" json:"books"`
}
