package shared

import (
	"github.com/golang-jwt/jwt/v5"
	"xxx/real_time/models"
)

// UserTokenClaims represents the structure of the user's ephemeral token
type UserTokenClaims struct {
	UserId    string
	SessionId string
	Role      models.Role
	jwt.RegisteredClaims
}
