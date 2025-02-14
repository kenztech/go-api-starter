package utils

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var (
	validate  *validator.Validate
	secretKey []byte
	otpStore  = make(map[string]string)
	otpMutex  = sync.Mutex{}
)

func init() {
	validate = validator.New()
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}
	jwtKey := os.Getenv("JWT_KEY")
	secretKey = []byte(jwtKey)
}

// StoreOTP stores the OTP for the given email.
func StoreOTP(email, otp string) {
	otpMutex.Lock()
	defer otpMutex.Unlock()

	otpStore[email] = otp

	go func() {
		time.Sleep(10 * time.Minute)
		otpMutex.Lock()
		delete(otpStore, email)
		otpMutex.Unlock()
	}()
}

func RetrieveOTP(email string) (string, bool) {
	otpMutex.Lock()
	defer otpMutex.Unlock()

	otp, exists := otpStore[email]
	return otp, exists
}

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func ValidateStruct(w http.ResponseWriter, request interface{}) bool {
	if err := validate.Struct(request); err != nil {
		var errs []string
		for _, err := range err.(validator.ValidationErrors) {
			field := strings.ToLower(err.Field())
			var message string
			switch err.Tag() {
			case "required":
				message = field + " is required"
			case "email":
				message = field + " must be a valid email address"
			case "min":
				message = field + " must be at least " + err.Param() + " characters long"
			default:
				message = field + " is not valid"
			}

			errs = append(errs, message)
		}

		SendError(w, strings.Join(errs, ", "), http.StatusBadRequest)
		return false
	}
	return true
}

func GenerateToken(email string, expiration time.Duration) (string, error) {
	expirationTime := time.Now().Add(expiration)

	claims := &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateToken(tokenString string) (string, bool) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return "", false
	}

	return claims.Email, true
}

func GenerateTokenWithOTP(email, otp string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"otp":   otp,
		"exp":   time.Now().Add(10 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func ValidateOTPToken(tokenStr, otp string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["otp"] != otp {
			return "", errors.New("invalid OTP")
		}

		if exp, ok := claims["exp"].(float64); ok && time.Unix(int64(exp), 0).Before(time.Now()) {
			return "", errors.New("token expired")
		}

		return claims["email"].(string), nil
	}

	return "", errors.New("invalid token")
}

func GenerateResetToken(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(15 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func ValidateResetToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok && time.Unix(int64(exp), 0).Before(time.Now()) {
			return "", errors.New("token expired")
		}

		return claims["email"].(string), nil
	}

	return "", errors.New("invalid token")
}

func GenerateOTP() string {
	const otpLength = 6
	const otpChars = "0123456789"

	otp := make([]byte, otpLength)
	for i := 0; i < otpLength; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(otpChars))))
		if err != nil {
			panic(err)
		}
		otp[i] = otpChars[num.Int64()]
	}

	return string(otp)
}
