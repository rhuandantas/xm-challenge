package http

import (
	"context"
	"fmt"
	"github.com/joomcode/errorx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/rhuandantas/xm-challenge/config"
	"github.com/rhuandantas/xm-challenge/internal/adapters/http/handlers"
)

type Server struct {
	host           string
	Server         *echo.Echo
	companyHandler *handlers.Company
	configs        *config.Config
}

// NewAPIServer creates the main http with all configurations necessary
func NewAPIServer(companyHandler *handlers.Company, configs *config.Config) *Server {
	host := configs.Server.Port
	app := echo.New()

	app.HideBanner = true
	app.HidePort = true

	app.Pre(middleware.RemoveTrailingSlash())
	app.Use(middleware.Recover())
	app.Use(middleware.CORS())
	app.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			log.Info(fmt.Sprintf("URI: %s, status: %d", v.URI, v.Status))
			return nil
		},
	}))

	return &Server{
		host:           host,
		Server:         app,
		companyHandler: companyHandler,
	}
}

func (hs *Server) RegisterHandlers() {
	hs.companyHandler.RegisterRoutes(hs.Server)
}

// Start starts an application on specific port
func (hs *Server) Start() {
	hs.RegisterHandlers()
	ctx := context.Background()
	log.Info(ctx, fmt.Sprintf("Starting a http at http://%s", hs.host))
	err := hs.Server.Start(hs.host)
	if err != nil {
		log.Error(ctx, errorx.Decorate(err, "failed to start the http server"))
		return
	}
}
