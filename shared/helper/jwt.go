package helper

import (
	"errors"
	config "rent-application/configs"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	Level string `json:"level"`
	Email string `json:"email"`
	Id    string `json:"id"`
	jwt.StandardClaims
}

func GenerateJwt(issuer, level, email string) (string, error) {
	session, _ := strconv.Atoi(config.Cfg.Jwt.SessionLogin)
	claims := &Claims{
		Level: level,
		Email: email,
		StandardClaims: jwt.StandardClaims{
			Issuer:    issuer,
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(session) * 30).Unix(),
		},
	}

	tokens := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return tokens.SignedString([]byte(config.Cfg.Jwt.SecretKey))
}

func ParseJwt(cookie string) (string, string, error) {
	var claims Claims
	token, err := jwt.ParseWithClaims(cookie, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Cfg.Jwt.SecretKey), nil
	})
	if err != nil || !token.Valid {
		return "", "", err
	}

	return claims.Level, claims.StandardClaims.Issuer, nil
}

func ParseJwtBearer(authHeader string) (string, string, error) {
	if authHeader == "" {
		return "", "", errors.New("authorization header is missing")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", "", errors.New("invalid authorization header format")
	}

	tokenString := parts[1]

	var claims Claims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Cfg.Jwt.SecretKey), nil
	})
	if err != nil || !token.Valid {
		return "", "", err
	}

	return claims.Level, claims.StandardClaims.Issuer, nil
}
