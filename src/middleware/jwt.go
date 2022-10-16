package middleware

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/idprm/go-yellowclinic/src/config"
	"github.com/idprm/go-yellowclinic/src/model"
)

func GenerateJWTToken(user model.User) (string, int64, error) {
	exp := time.Now().Add(time.Hour * 24).Unix()
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["exp"] = exp

	t, err := token.SignedString([]byte(config.ViperEnv("JWT_SECRET_AUTH")))

	if err != nil {
		return "", 0, err
	}

	return t, exp, nil
}

func RefreshJWTToken() error {
	return nil
}
