package service

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

const (
	signingKey = "qrkjk#4#%35FSFJlja#4353KSFjH"
	tokenTTL   = 12 * time.Hour
)

type Claims struct {
	UserID   string `json:"user_id"`
	ClientIP string `json:"client_ip"`
	jwt.StandardClaims
}

func GenerateAccessToken(guid string, clientIp string) (string, error) {
	claims := &Claims{
		UserID:   guid,
		ClientIP: clientIp,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(signingKey))
}

func ValidateAccessToken(token string) (*Claims, error) {
	tok, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err!=nil{
		return nil, err
	}

	if claims, ok := tok.Claims.(*Claims); ok && tok.Valid {
		return claims, nil
	}
	return nil, err
}

func GenerateRefreshToken() string {
	tokenBytes := sha256.Sum256([]byte(time.Now().String()))
	return base64.StdEncoding.EncodeToString(tokenBytes[:])
}

func HashRefreshToken(token string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(token),bcrypt.DefaultCost)
}

func CompareRefreshToken(hashedToken, token string) error {

	log.Print(hashedToken)
	bcrypt.GenerateFromPassword([]byte(token),bcrypt.DefaultCost)

	return bcrypt.CompareHashAndPassword([]byte(hashedToken), []byte(token))
}
