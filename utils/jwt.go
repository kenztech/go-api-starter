package utils

import (
	"context"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type contextKey string

const UserContextKey contextKey = "userData"

type UserData struct {
	ID    uint
	Email string
	Role  string
}

func SetUserDataInContext(ctx context.Context, email, role string) context.Context {
	userData := UserData{
		Email: email,
		Role:  role,
	}
	return context.WithValue(ctx, UserContextKey, userData)
}

func GetUserDataFromContext(ctx context.Context) (UserData, bool) {
	userData, ok := ctx.Value(UserContextKey).(UserData)
	return userData, ok
}

// GenerateAccessToken creates a JWT including user ID, email, and role
func GenerateAccessToken(id uint, email, role string) (string, error) {
	claims := jwt.MapClaims{
		"authorized": true,
		"id":         id,
		"email":      email,
		"role":       role,
		"exp":        time.Now().Add(15 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// GenerateRefreshToken creates a refresh token with user ID, email, and role
func GenerateRefreshToken(email, role string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"role":  role,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateTheToken parses and validates a JWT
func ValidateTheToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
}

// VerifyRefreshToken extracts user details from the refresh token
func VerifyRefreshToken(tokenString string) (string, string, error) {
	token, err := ValidateTheToken(tokenString)
	if err != nil || !token.Valid {
		return "", "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", err
	}

	email, ok := claims["email"].(string)
	if !ok {
		return "", "", err
	}

	role, ok := claims["role"].(string)
	if !ok {
		return "", "", err
	}

	return email, role, nil
}
