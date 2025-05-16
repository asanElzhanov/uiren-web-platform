package jwt_maker

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtMaker struct {
	duration time.Duration
}

func NewJWTMaker(duration time.Duration) *jwtMaker {
	return &jwtMaker{
		duration: duration,
	}
}

func (maker *jwtMaker) NewToken(payload PayloadDTO) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	payload.Duration = maker.duration

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = payload.ID
	claims["email"] = payload.Email
	claims["username"] = payload.Username
	claims["firstname"] = payload.Firstname
	claims["lastname"] = payload.Lastname
	claims["isAdmin"] = payload.IsAdmin
	claims["exp"] = time.Now().Add(payload.Duration).Unix()

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
