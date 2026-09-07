package models

type Book struct {
	Model
	Isbn          string `gorm:"not null" json:"isbn"`
	Title         string `gorm:"not null" json:"title"`
	Author        string `gorm:"not null" json:"author"`
	Description   string `json:"description"`
	ThumbnailLink string `json:"thumbnail_link"`
	Pages         int    `json:"pages"`
	Owned         bool   `gorm:"default:false" json:"owned"`
	Read          bool   `gorm:"default:false" json:"read"`

	UserId uint `gorm:"not null"` // Foreign key to assign books to a user
}

type LookupBook struct {
	Isbn          string
	Title         string
	Author        string
	Description   string
	ThumbnailLink string
	Pages         int
}
