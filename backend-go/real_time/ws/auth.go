package ws

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"xxx/real_time/config"
	"xxx/shared"
)

func extractTokenData(tokenString string) (*shared.UserToken, error) {
	token, err := jwt.ParseWithClaims(tokenString, &shared.UserToken{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.LoadConfig().JWT.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	switch {
	case token.Valid:
		fmt.Println("OK token")
	case errors.Is(err, jwt.ErrTokenMalformed):
		fmt.Println("Malformed token")
	default:
		fmt.Println("Couldn't handle this token:", err)
	}

	claims, ok := token.Claims.(*shared.UserToken)
	fmt.Println("Decoded token: ", *claims)
	if !ok {
		return nil, fmt.Errorf("error decoding jwt")
	}

	return claims, nil
}

// TODO: func validateToken() {}
