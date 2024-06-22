package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"go-template/pkg/common/env"
	"time"
)

func SignToken(userId uint, expireTime time.Time) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = expireTime.Unix()
	claims["userId"] = userId

	tokenString, err := token.SignedString([]byte(env.JWT_SECRET))
	if err != nil {
		return ""
	}

	return tokenString
}
