/*
 * @Author       : jayj
 * @Date         : 2021-06-24 09:40:59
 * @Description  :
 */
package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	salt   string        = "etcdgate"
	expire time.Duration = time.Duration(time.Hour * 168) // 7 days
)

var jwtSecret = []byte(salt)

type Claims struct {
	Address  string `json:"address"`
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

// GenerateToken
func GenerateToken(address, username, password string) (token string, err error) {

	now := time.Now()
	expireTime := now.Add(expire)

	claims := Claims{
		address,
		username,
		password,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "gin-blog",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = tokenClaims.SignedString(jwtSecret)

	return token, err
}

// ParseToken get claims
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		} else if ok {
			return claims, err
		}
	}

	return nil, err
}
