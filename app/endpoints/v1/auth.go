package v1

import (
	"ecommerce/app/core"
	"ecommerce/app/core/middlewares"
	"ecommerce/app/crud"
	"ecommerce/app/models"
	"ecommerce/app/schemas"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

// Login
// @Summary Authenticate user
// @Description Authenticates a user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body schemas.UserLogin true "User login credentials"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]interface{}
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var DB = core.GetDB()

	var request schemas.UserLogin
	if err := c.ShouldBindJSON(&request); err != nil {
		core.HandleValidationErrors(c, err)
		return
	}
	accessToken, refreshToken, err := crud.LoginUser(DB, request.Email, request.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"accessToken": accessToken, "refreshToken": refreshToken})
}

// Register
// @Summary Register a new user
// @Description Creates a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body schemas.UserRegister true "User registration details"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /auth/register [post]
func Register(c *gin.Context) {
	var DB = core.GetDB()
	var request schemas.UserRegister

	// check the request data
	if err := c.ShouldBindJSON(&request); err != nil {
		core.HandleValidationErrors(c, err)
		return
	}
	// create new accessToken / refreshToken
	accessToken, refreshToken, err := crud.CreateUser(DB, request.Email, request.PhoneNumber, request.Password, request.FirstName, request.LastName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"accessToken": accessToken, "refreshToken": refreshToken})
}

// Refresh
// @Summary Refresh access token
// @Description Generates a new access token using the refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body schemas.RefreshToken true "Refresh token request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]interface{}
// @Router /auth/refresh [post]
func Refresh(c *gin.Context) {
	var DB = core.GetDB()
	var request schemas.RefreshToken

	// check the request data
	if err := c.ShouldBindJSON(&request); err != nil {
		core.HandleValidationErrors(c, err)
		return
	}
	// create new access token
	accessToken, err := crud.RefreshUserToken(DB, request.RefreshToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// return the new access token
	c.JSON(http.StatusOK, gin.H{"accessToken": accessToken})
}

// Logout
// @Summary Logout user
// @Description Invalidates the user's access and refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /auth/logout [post]
func Logout(c *gin.Context) {
	var DB = core.GetDB()
	var request schemas.RefreshToken
	if err := c.ShouldBindJSON(&request); err != nil {
		core.HandleValidationErrors(c, err)
		return
	}
	if err := crud.BlackListRefreshToken(DB, request.RefreshToken); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User logged out successfully!"})
}

// ChangePassword
// @Summary Change user's password
// @Description Allows an authenticated user to change their password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body schemas.UserChangePasswordSchema true "Change password request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /auth/change-password [post]
func ChangePassword(c *gin.Context) {
	var request schemas.UserChangePasswordSchema
	if err := c.ShouldBindJSON(&request); err != nil {
		core.HandleValidationErrors(c, err)
		return
	}
	user := c.MustGet("user").(models.User)
	fmt.Println(user.Email)
	DB := core.GetDB()
	if err := crud.UpdateUserPassword(DB, user, request.NewPassword, request.OldPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User changed password successfully!"})
}

// VerifyEmail
// @Summary Verify user's email
// @Description Verifies the user's email address using a token
// @Tags auth
// @Accept json
// @Produce json
// @Param token path string true "Email verification token"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /auth/verify/{token} [get]
func VerifyEmail(c *gin.Context) {
	var DB = core.GetDB()

	tokenStr := c.Param("token")
	token, err := uuid.Parse(tokenStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := crud.VerifyUser(DB, token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User verified successfully"})
}

// ResendVerificationEmail
// @Summary Resend email verification
// @Description Resends the email verification link to the user
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /auth/resend-verify [post]
func ResendVerificationEmail(c *gin.Context) {
	var DB = core.GetDB()
	user := c.MustGet("user").(models.User)

	if err := crud.ResendVerifyUser(DB, user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User resending email successfully!"})
}

// ResetPasswordRequest
// @Summary Request password reset
// @Description Initiates the password reset process for a user
// @Tags auth
// @Accept json
// @Produce json
// @Param email query string true "User's email address"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /auth/reset-password [post]
func ResetPasswordRequest(c *gin.Context) {
	type requestSchema struct {
		Email string `json:"email" binding:"required"`
	}
	var request requestSchema
	DB := core.GetDB()
	if err := c.ShouldBindJSON(&request); err != nil {
		core.HandleValidationErrors(c, err)
		return
	}
	if err := crud.SendPasswordResetToken(DB, request.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User reset password has been sent successfully!"})
}

// ResetPassword
// @Summary Reset user's password
// @Description Resets the user's password using a token
// @Tags auth
// @Accept json
// @Produce json
// @Param token path string true "Password reset token"
// @Param new_password query string true "User's New Password"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /auth/reset-password/{token} [post]
func ResetPassword(c *gin.Context) {
	type requestSchema struct {
		NewPassword string `json:"new_password" binding:"required"`
	}
	var request requestSchema

	if err := c.ShouldBindJSON(&request); err != nil {
		core.HandleValidationErrors(c, err)
		return
	}

	var DB = core.GetDB()
	tokenStr := c.Param("token")
	token, err := uuid.Parse(tokenStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := crud.ResetUserPassword(DB, token, request.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User password has been reset successfully!"})
}

// UpdateUser
// @Summary Update user details
// @Description Allows an authenticated user to update their profile details
// @Tags auth
// @Accept json
// @Produce json
// @Param request body schemas.UpdateUserRequest true "Update user request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /auth/update-user [patch]
func UpdateUser(c *gin.Context) {
	DB := core.GetDB()
	user := c.MustGet("user").(models.User)

	var request schemas.UpdateUserRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		core.HandleValidationErrors(c, err)
		return
	}

	if err := crud.UpdateUserInfo(DB, user.ID, request.FirstName, request.LastName, request.Email, request.PhoneNumber); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully!"})

}

func AuthRouter(router *gin.Engine) {
	// Define the Router
	public := router.Group("/api/v1/auth")
	{
		public.POST("/login", Login)
		public.POST("/register", Register)
		public.POST("/refresh", Refresh)
		public.POST("/logout", Logout)
		public.GET("/verify/:token", VerifyEmail)
		public.POST("/reset-password", ResetPasswordRequest)
		public.POST("/reset-password/:token", ResetPassword)
	}
	protected := router.Group("/api/v1/auth")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.POST("/change-password", ChangePassword)
		protected.POST("/resend-verify", ResendVerificationEmail)
		protected.PATCH("/update-user", UpdateUser)
	}
}
