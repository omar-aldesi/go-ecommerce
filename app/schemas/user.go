package schemas

type UserBase struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=32"`
}

type UserLogin struct {
	UserBase
}

type UserRegister struct {
	UserBase
	PhoneNumber string `json:"phone_number" binding:"required,e164"`
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
}

type RefreshToken struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
type UserChangePasswordSchema struct {
	OldPassword string `json:"old_password" binding:"required,min=8,max=32"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=32"`
}

type UpdateUserRequest struct {
	FirstName   *string `json:"first_name"`
	LastName    *string `json:"last_name"`
	Email       *string `json:"email"`
	PhoneNumber *string `json:"phone_number"`
}
