package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"time"
	"user-service/errors"
)

const AccessTokenValidity = time.Hour * 24
const RefreshTokenValidity = time.Hour * 24

// verifyAccessToken verifies a token
func verifyToken(tokenString string, secret string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
}

func isJWTSecretEmpty(secret string) bool {
	return secret == ""
}

func isAccessTokenEmpty(token string) bool {
	return token == ""
}

func ValidateToken(token string, secret string) (*jwt.Token, error) {
	tk, err := verifyToken(token, secret)
	if err != nil {
		log.Println(err)                                 // TODO: remove
		return nil, fmt.Errorf("invalid token: %v", err) // TODO: probably need to errors.NEw
	}
	if !tk.Valid {
		return nil, errors.New("invalid token", http.StatusUnauthorized)
	}
	return tk, nil
}

func getClaims(token *jwt.Token) (jwt.MapClaims, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("could not get claims")
	}

	return claims, claims.Valid()
}

func ValidateAndGetClaims(tokenString string, secret string) (jwt.MapClaims, error) {
	if tokenString == "" {
		return nil, errors.New("invalid token (token is empty)", http.StatusUnauthorized)
	}
	token, err := ValidateToken(tokenString, secret)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %v", err)
	}
	claims, err := getClaims(token)
	if err != nil {
		return nil, fmt.Errorf("failed to get claims: %v", err)
	}

	return claims, nil
}

type Claim struct {
	UserId string `json:"user_id"`
	Email  string `json:"email"`
}

// GenerateToken generates only an access token
func GenerateToken(email string, id int, secret string) (string, error) {
	if secret == "" {
		return "", errors.New("", http.StatusInternalServerError)
	}
	// Generate claims
	claims := GenerateClaims(email, id)
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateClaims(email string, id int) jwt.MapClaims {
	accessClaims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(AccessTokenValidity).Unix(),
		"id":    id,
	}
	return accessClaims
}
