package security

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"time"
)

var (
	accessSecret  = []byte("your_access_secret_key")
	refreshSecret = []byte("your_refresh_secret_key")
)

var (
	accessExpiryMinutes time.Duration = 1000
	refreshExpiryDays   time.Duration = 7
)

// Claims structure
type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// createToken creates a token with given claims and secret
func createToken(claims Claims, secret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func CreateAccessToken(email string) (string, error) {
	// Create the access token
	accessClaims := Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessExpiryMinutes * time.Minute)),
		},
	}
	accessToken, err := createToken(accessClaims, accessSecret)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func CreateRefreshToken(email string) (string, error) {
	// Create the refresh token
	refreshClaims := Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshExpiryDays * 24 * time.Hour)),
		},
	}
	refreshToken, err := createToken(refreshClaims, refreshSecret)
	if err != nil {
		return "", err
	}
	return refreshToken, nil
}

func CreateJwtDefaultTokens(email string) (string, string, error) {
	accessToken, err := CreateAccessToken(email)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := CreateRefreshToken(email)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

// ValidateToken validates the token string and returns the claims if valid
func ValidateToken(db *gorm.DB, tokenString, tokenType string) (*Claims, error) {
	if tokenType == "refresh" && IsTokenBlacklisted(db, tokenString) {
		return nil, fmt.Errorf("token is blacklisted")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		if tokenType == "access" {
			return accessSecret, nil
		} else {
			return refreshSecret, nil
		}
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
