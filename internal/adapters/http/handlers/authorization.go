package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/rhuandantas/xm-challenge/internal/adapters/http/middlewares/auth"
)

type Authorization struct {
	jwt auth.Token
}

func NewAuthorizationHandler(token auth.Token) *Authorization {
	return &Authorization{
		jwt: token,
	}
}

func (a *Authorization) RegisterRoute(server *echo.Echo) {
	server.GET("/auth", a.GetToken)

}

func (a *Authorization) GetToken(ctx echo.Context) error {
	token, err := a.jwt.GenerateToken()
	if err != nil {
		return ctx.JSON(500, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(200, map[string]string{"token": token})
}
