package shared

import (
	"github.com/golang-jwt/jwt/v5"
	"xxx/real_time/models"
)

const JwtKey = "some_complex_key"

type UserTokenClaims struct {
	UserId    string
	SessionId string
	Role      models.Role
	jwt.RegisteredClaims
}
