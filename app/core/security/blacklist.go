package security

import (
	"ecommerce/app/models"
	"gorm.io/gorm"
	"time"
)

func BlacklistToken(db *gorm.DB, token string, expiresAt time.Time) error {
	blacklistedToken := models.BlacklistedToken{
		Token:     token,
		ExpiresAt: expiresAt,
	}
	return db.Create(&blacklistedToken).Error
}

func IsTokenBlacklisted(db *gorm.DB, token string) bool {
	var count int64
	db.Model(&models.BlacklistedToken{}).Where("token = ? AND expires_at > ?", token, time.Now()).Count(&count)
	return count > 0
}

func CleanupBlacklist(db *gorm.DB) error {
	return db.Where("expires_at < ?", time.Now()).Delete(&models.BlacklistedToken{}).Error
}
