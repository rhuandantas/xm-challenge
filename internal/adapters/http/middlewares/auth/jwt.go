package auth

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/rhuandantas/fiap-tech-challenge-commons/pkg/errors"
	"github.com/rhuandantas/xm-challenge/config"
	"strings"
	"time"
)

type Token interface {
	GenerateToken() (string, error)
	VerifyToken(next echo.HandlerFunc) echo.HandlerFunc
}

type JwtToken struct {
	cfg *config.Config
}

func NewJwtToken(cfg *config.Config) Token {
	return &JwtToken{
		cfg: cfg,
	}
}

type jwtCustomClaims struct {
	jwt.RegisteredClaims
}

func (jt *JwtToken) GenerateToken() (string, error) {
	secret := []byte(jt.cfg.JWT.Secret)
	claims := &jwtCustomClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jt.cfg.JWT.Expiration)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return t, nil
}

func (jt *JwtToken) VerifyToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenStr := jt.getToken(c)
		if tokenStr == "" {
			return errors.HandleError(c, errors.Unauthorized.New("authentication key not found"))
		}

		secret := []byte(jt.cfg.JWT.Secret)
		tkn, err := jwt.ParseWithClaims(tokenStr, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})
		if err != nil {
			return errors.HandleError(c, errors.Unauthorized.New(err.Error()))
		}

		if !tkn.Valid {
			return errors.HandleError(c, errors.Unauthorized.New("authentication is not valid"))
		}

		return next(c)
	}
}

func (jt *JwtToken) getToken(c echo.Context) string {
	tokenStr := ""
	if bearer := c.Request().Header.Get("Authorization"); bearer != "" {
		if strings.Contains(bearer, "Bearer") {
			tokenStr = strings.Split(bearer, " ")[1]
		}
	}

	if tokenStr == "" {
		tokenStr = c.Request().Header.Get("token")
	}

	return tokenStr
}
