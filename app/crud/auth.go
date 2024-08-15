package crud

import (
	"ecommerce/app/core"
	"ecommerce/app/core/security"
	"ecommerce/app/models"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

func LoginUser(db *gorm.DB, email, password string) (string, string, error) {
	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return "", "", err
	}

	if !security.CheckPasswordHash(password, user.HashedPassword) {
		return "", "", fmt.Errorf("password is not correct")
	}
	accessToken, refreshToken, err := security.CreateJwtDefaultTokens(email)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func CreateUser(db *gorm.DB, email, phoneNumber, password, firstName, lastName string) (string, string, error) {
	var user models.User
	// check if the Email and the phoneNumber are unique
	if err := db.Where("email = ?", email).Or("phone_number = ?", phoneNumber).First(&user).Error; err == nil {
		return "", "", fmt.Errorf("user already exists")
	}
	user = models.User{
		FirstName:   firstName,
		LastName:    lastName,
		Email:       email,
		PhoneNumber: phoneNumber,
	}
	// hash and set the password
	hashedPassword, err := security.HashPassword(password)
	if err != nil {
		return "", "", err
	}
	user.HashedPassword = hashedPassword

	// create new user in the db
	if err := db.Create(&user).Error; err != nil {
		return "", "", err
	}
	// create new accessToken / refreshToken
	accessToken, refreshToken, err := security.CreateJwtDefaultTokens(email)
	if err != nil {
		return "", "", err
	}
	defer db.Create(&models.EmailVerificationToken{
		UUID:   uuid.New(),
		UserID: user.ID,
	})
	return accessToken, refreshToken, nil
}

func RefreshUserToken(db *gorm.DB, refreshToken string) (string, error) {
	claims, err := security.ValidateToken(db, refreshToken, "refresh")
	if err != nil {
		return "", err
	}

	accessToken, err := security.CreateAccessToken(claims.Email)
	if err != nil {
		return "", err
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return "", err
	}

	if err := UpdateUserLastLogin(tx, claims.Email); err != nil {
		tx.Rollback()
		return "", err
	}

	if err := tx.Commit().Error; err != nil {
		return "", err
	}

	return accessToken, nil
}

func UpdateUserLastLogin(db *gorm.DB, email string) error {
	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return err
	}
	user.LastLogin = time.Now()
	if err := db.Save(&user).Error; err != nil {
		return err
	}
	return nil
}

func BlackListRefreshToken(db *gorm.DB, refreshToken string) error {
	claims, err := security.ValidateToken(db, refreshToken, "refresh")
	if err != nil {
		return err
	}
	return security.BlacklistToken(db, refreshToken, claims.ExpiresAt.Time)
}

func UpdateUserPassword(db *gorm.DB, user models.User, newPassword, oldPassword string) error {
	// check if the old_password == user current password
	if !security.CheckPasswordHash(oldPassword, user.HashedPassword) {
		return fmt.Errorf("password is not correct")
	}
	// hash the new password
	hashedPassword, err := security.HashPassword(newPassword)
	if err != nil {
		return err
	}
	// save the new user password
	user.HashedPassword = hashedPassword
	db.Save(&user)
	return nil
}

func SendPasswordResetToken(db *gorm.DB, email string) error {
	var user models.User
	if err := db.First(&user, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user not found")
		}
		return err
	}
	if err := db.Where("user_id = ?", user.ID).Delete(&models.PasswordResetToken{}).Error; err != nil {
		return fmt.Errorf("error deleting old tokens: %w", err)
	}
	passwordResetToken := models.PasswordResetToken{
		UUID:      uuid.New(),
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	if err := db.Create(&passwordResetToken).Error; err != nil {
		return fmt.Errorf("error creating token: %w", err)
	}
	if err := core.SendEmail(user.Email, "Password Reset Token", passwordResetToken.UUID); err != nil {
		return err
	}
	return nil
}

func ResetUserPassword(db *gorm.DB, token uuid.UUID, newPassword string) error {
	var dbToken models.PasswordResetToken
	// check the Token if exists in the db
	if err := db.First(&dbToken, "uuid = ?", token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("invalid or expired token")
		}
		return fmt.Errorf("error finding token: %w", err)
	}
	// hash the new password
	hashedPassword, err := security.HashPassword(newPassword)
	if err != nil {
		return err
	}
	// Update the user's password
	if err := db.Model(&models.User{}).Where("id = ?", dbToken.UserID).Update("hashed_password", hashedPassword).Error; err != nil {
		return fmt.Errorf("error updating password: %w", err)
	}
	// delete all user tokens from the db
	if err := db.Where("user_id = ?", dbToken.UserID).Delete(&models.PasswordResetToken{}).Error; err != nil {
		return fmt.Errorf("error deleting token: %w", err)
	}
	return nil
}

func VerifyUser(db *gorm.DB, token uuid.UUID) error {

	var dbToken models.EmailVerificationToken

	// Find the token and preload the associated user
	if err := db.First(&dbToken, "uuid = ?", token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("invalid or expired token")
		}
		return fmt.Errorf("error finding token: %w", err)
	}

	// Check if user is already verified
	if dbToken.User.IsVerified {
		return fmt.Errorf("user already verified")
	}

	// Delete all verification tokens for this user
	if err := db.Where("user_id = ?", dbToken.UserID).Delete(&models.EmailVerificationToken{}).Error; err != nil {
		return fmt.Errorf("error deleting old tokens: %w", err)
	}

	// verify the user
	if err := db.Model(&models.User{}).Where("id = ?", dbToken.UserID).Update("is_verified", true).Error; err != nil {
		return err
	}

	return nil
}

func ResendVerifyUser(db *gorm.DB, user models.User) error {
	if user.IsVerified {
		return fmt.Errorf("user already verfied")
	}
	// deleting all old tokens
	if err := db.Where("user_id = ?", user.ID).Delete(&models.EmailVerificationToken{}).Error; err != nil {
		return fmt.Errorf("error deleting old tokens: %w", err)
	}
	// create new token
	token := models.EmailVerificationToken{
		UUID:      uuid.New(),
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	if err := db.Create(&token).Error; err != nil {
		return fmt.Errorf("error creating token: %w", err)
	}
	// send the new token via email
	if err := core.SendEmail(user.Email, "Email Verification Token", token.UUID); err != nil {
		return err
	}
	return nil
}

func UpdateUserInfo(db *gorm.DB, userID uint, firstName, lastName, email, phoneNumber *string) error {
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return err
	}

	// Debug: Print user before update
	fmt.Printf("User before update: %+v\n", user)

	updates := make(map[string]interface{})
	if firstName != nil {
		updates["first_name"] = *firstName
	}
	if lastName != nil {
		updates["last_name"] = *lastName
	}
	if email != nil {
		updates["email"] = *email
	}
	if phoneNumber != nil {
		updates["phone_number"] = *phoneNumber
	}

	if len(updates) > 0 {
		if err := db.Model(&user).Updates(updates).Error; err != nil {
			return err
		}
	}
	// Refresh user data
	if err := db.First(&user, userID).Error; err != nil {
		return err
	}
	return nil
}
