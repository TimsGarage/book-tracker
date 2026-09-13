package models

type Book struct {
	Model
	Isbn          string `gorm:"not null" json:"isbn" form:"isbn"`
	Title         string `gorm:"not null" json:"title" form:"title"`
	Author        string `gorm:"not null" json:"author" form:"author"`
	Description   string `json:"description" form:"description"`
	Release       string `json:"release" form:"release"`
	Publisher     string `json:"publisher" form:"publisher"`
	ThumbnailLink string `json:"thumbnail_link" form:"thumbnail_link"`
	Pages         int    `json:"pages" form:"pages"`

	OwnershipStatus string `gorm:"type:varchar(20);check:ownership_status IN ('wishlist', 'unowned', 'owned', 'borrowed')" json:"ownership_status" form:"ownership_status"`
	OwnedSince      string `json:"owned_since" form:"owned_since"`
	ReadingStatus   string `gorm:"type:varchar(20);check:reading_status IN ('unread', 'reading', 'read')" json:"reading_status" form:"reading_status"`

	UserId uint `gorm:"not null"` // Foreign key to assign books to a user
}

type LookupBook struct {
	Isbn          string `json:"isbn"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	Description   string `json:"description"`
	Release       string `json:"release"`
	Publisher     string `json:"publisher"`
	ThumbnailLink string `json:"thumbnail_link"`
	Pages         int    `json:"pages"`
}
