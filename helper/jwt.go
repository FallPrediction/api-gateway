package helper

import (
	"api-gateway/logger"
	"errors"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type Jwt struct{}

type UserData struct {
	UserId uint32 `json:"user_id"`
	Email  string `json:"email"`
	Scope  string `json:"scope"`
}

type UserClaims struct {
	UserData
	jwt.RegisteredClaims
}

func (c UserClaims) Validate() error {
	if c.ExpiresAt == nil {
		return errors.New("Claim expires at should exists.")
	}
	return nil
}

func (j *Jwt) ParseUserToken(tokenString string) (*UserClaims, error) {
	logger := logger.NewLogger()
	key := os.Getenv("APP_KEY")
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		logger.Error("Parse token failed.", zap.String("err", fmt.Sprintf("%v", err)))
		return nil, err
	} else if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("the provided token is invalid")
}
