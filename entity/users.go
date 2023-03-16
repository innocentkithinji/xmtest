package entity

import "github.com/golang-jwt/jwt/v4"

type (
	User struct {
		ID       string `json:"id" bson:"_id"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	LoginCredentials struct {
		Token string `json:"token"`
	}

	JWTClaims struct {
		Email string
		UID   string
		jwt.RegisteredClaims
	}
)
