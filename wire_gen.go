// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/rhuandantas/xm-challenge/config"
	"github.com/rhuandantas/xm-challenge/internal/adapters/async/kafka"
	"github.com/rhuandantas/xm-challenge/internal/adapters/http"
	"github.com/rhuandantas/xm-challenge/internal/adapters/http/handlers"
	"github.com/rhuandantas/xm-challenge/internal/adapters/http/middlewares"
	"github.com/rhuandantas/xm-challenge/internal/adapters/http/middlewares/auth"
	"github.com/rhuandantas/xm-challenge/internal/adapters/repository"
	"github.com/rhuandantas/xm-challenge/internal/adapters/repository/mysql"
	"github.com/rhuandantas/xm-challenge/internal/core/usecases"
)

// Injectors from wire.go:

func InitializeWebServer() (*http.Server, error) {
	configConfig := config.LoadConfig()
	dbConnector := mysql.NewMySQLConnector(configConfig)
	companyRepo := repository.NewCompanyRepo(dbConnector)
	getCompany := usecases.NewGetCompany(companyRepo)
	producer, err := kafka.NewProducer(configConfig)
	if err != nil {
		return nil, err
	}
	createCompany := usecases.NewCreateCompany(companyRepo, producer)
	deleteCompany := usecases.NewDeleteCompany(companyRepo, producer)
	updateCompany := usecases.NewUpdateCompany(companyRepo, producer)
	validator := middlewares.NewCustomValidator()
	token := auth.NewJwtToken(configConfig)
	company := handlers.NewCompanyHandler(getCompany, createCompany, deleteCompany, updateCompany, validator, token)
	authorization := handlers.NewAuthorizationHandler(token)
	server := http.NewAPIServer(company, authorization, configConfig)
	return server, nil
}
