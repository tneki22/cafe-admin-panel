package vars

import "github.com/golang-jwt/jwt/v5"

type JWTClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
