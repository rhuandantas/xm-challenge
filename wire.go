// wire.go
//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/rhuandantas/xm-challenge/config"
	"github.com/rhuandantas/xm-challenge/internal/adapters/http"
	"github.com/rhuandantas/xm-challenge/internal/adapters/http/handlers"
	"github.com/rhuandantas/xm-challenge/internal/adapters/http/middlewares"
	"github.com/rhuandantas/xm-challenge/internal/adapters/http/middlewares/auth"
	"github.com/rhuandantas/xm-challenge/internal/adapters/messaging/kafka"
	"github.com/rhuandantas/xm-challenge/internal/adapters/repository"
	"github.com/rhuandantas/xm-challenge/internal/adapters/repository/mysql"
	"github.com/rhuandantas/xm-challenge/internal/core/usecases"
)

func InitializeWebServer() (*http.Server, error) {
	wire.Build(config.LoadConfig,
		auth.NewJwtToken,
		mysql.NewMySQLConnector,
		kafka.NewProducer,
		repository.NewCompanyRepo,
		usecases.NewCreateCompany,
		usecases.NewGetCompany,
		usecases.NewDeleteCompany,
		usecases.NewUpdateCompany,
		middlewares.NewCustomValidator,
		handlers.NewCompanyHandler,
		handlers.NewAuthorizationHandler,
		http.NewAPIServer)
	return &http.Server{}, nil
}
