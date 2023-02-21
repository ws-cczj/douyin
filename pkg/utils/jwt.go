package utils

import (
	"douyin/consts"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	mySecret        = []byte(consts.JWTSecret)
	ErrInvalidToken = errors.New("verify Token Failed")
)

type MyClaim struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

// GenToken 颁发token
func GenToken(UserID int64) (token string, err error) {
	mc := MyClaim{
		UserID: UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: getJWTTime(consts.JWTTokenExpiredAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    consts.JWTIssuer,
			Subject:   consts.JWTDouyin,
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, mc).SignedString(mySecret)
}

// VerifyToken 验证Token
func VerifyToken(tokenStr string) (*MyClaim, error) {
	var mc = new(MyClaim)
	token, err := jwt.ParseWithClaims(tokenStr, mc, func(token *jwt.Token) (interface{}, error) {
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, ErrInvalidToken
	}
	return mc, nil
}

func getJWTTime(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(time.Duration(t) * time.Second))
}
