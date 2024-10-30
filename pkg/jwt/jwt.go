package jwt

import (
	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/config"
	token "github.com/golang-jwt/jwt/v5"
)

type Jwt struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Jwt {
	return &Jwt{cfg: cfg}
}

type claims struct {
	token.RegisteredClaims
	UserID int
}

func (j *Jwt) GenerateTokenForUser(userID int) (string, error) {
	jwtToken := token.NewWithClaims(token.SigningMethodHS256, claims{
		RegisteredClaims: token.RegisteredClaims{},
		UserID:           userID,
	})

	jwtString, err := jwtToken.SignedString([]byte(j.cfg.Secret))
	return jwtString, err
}

func (j *Jwt) GetUserID(tokenString string) int {
	claims := &claims{}

	jwtToken, err := token.ParseWithClaims(tokenString, claims, func(t *token.Token) (interface{}, error) {
		return []byte(j.cfg.Secret), nil
	})
	if err != nil {
		return -1
	}

	if !jwtToken.Valid {
		return -1
	}

	return claims.UserID
}
