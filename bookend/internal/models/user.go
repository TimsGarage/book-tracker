package models

type User struct {
	Model
	Username string `gorm:"not null,unique" json:"username"`
	Password string `gorm:"not null" json:"-"`

	Books []Book `gorm:"foreignKey:UserId" json:"books"`
}

// RegisterInput defines the required payload for registering a new user
type RegisterInput struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginInput defines the credentials required for authentication
type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserResponse is the safe representation of a user without sensitive data
type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

// AuthResponse returns the generated token and safe user profile
type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// ChangePasswordInput defines the payload for updating a user's password
type ChangePasswordInput struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}
