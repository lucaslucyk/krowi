package security

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/lucaslucyk/krowi/pkg/config"
	"github.com/lucaslucyk/krowi/pkg/models"
)

func CreateToken(user *models.User, exp time.Duration) (string, error) {
	cfg, err := config.New()
	if err != nil {
		return "", err
	}
	tokenByte := jwt.New(jwt.SigningMethodHS256)
	now := time.Now().UTC()
	claims := tokenByte.Claims.(jwt.MapClaims)
	expDuration := exp

	claims["sub"] = user.ID
	claims["admin"] = user.IsAdmin
	claims["exp"] = now.Add(expDuration).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	return tokenByte.SignedString([]byte(cfg.SecretKey))
}
