package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/rhuandantas/xm-challenge/internal/adapters/http/middlewares"
	"github.com/rhuandantas/xm-challenge/internal/adapters/http/middlewares/auth"
	"github.com/rhuandantas/xm-challenge/internal/core/domain"
	"github.com/rhuandantas/xm-challenge/internal/core/usecases"
	"github.com/rhuandantas/xm-challenge/internal/errors"
	"net/http"
)

type Company struct {
	getCompanyUseCase    usecases.GetCompany
	createCompanyUseCase usecases.CreateCompany
	deleteCompanyUseCase usecases.DeleteCompany
	updateCompanyUseCase usecases.UpdateCompany
	validator            middlewares.Validator
	jwtMiddleware        auth.Token
}

func NewCompanyHandler(getCompanyUseCase usecases.GetCompany, createCompanyUseCase usecases.CreateCompany, deleteCompanyUseCase usecases.DeleteCompany, updateCompany usecases.UpdateCompany,
	validator middlewares.Validator, jwtMiddleware auth.Token) *Company {
	return &Company{
		getCompanyUseCase:    getCompanyUseCase,
		createCompanyUseCase: createCompanyUseCase,
		deleteCompanyUseCase: deleteCompanyUseCase,
		updateCompanyUseCase: updateCompany,
		validator:            validator,
		jwtMiddleware:        jwtMiddleware,
	}
}

func (p *Company) RegisterRoutes(server *echo.Echo) {
	server.GET("/companies/:name", p.getCompany)
	server.POST("/companies", p.storeCompany, p.jwtMiddleware.VerifyToken)
	server.DELETE("/companies/:id", p.deleteCompany, p.jwtMiddleware.VerifyToken)
	server.PATCH("/companies/:id", p.updateCompany, p.jwtMiddleware.VerifyToken)
}

func (p *Company) getCompany(ctx echo.Context) error {
	name := ctx.Param("name")
	if name == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "please send a company name"})
	}

	company, err := p.getCompanyUseCase.Execute(ctx.Request().Context(), name)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, company)
}

func (p *Company) storeCompany(ctx echo.Context) error {
	var company domain.Company
	if err := ctx.Bind(&company); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := company.Validate(); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := p.validateCompany(&company); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	_, err := p.createCompanyUseCase.Execute(ctx.Request().Context(), &company)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusCreated, map[string]string{"message": "company created"})
}

func (p *Company) validateCompany(company *domain.Company) error {
	if err := p.validator.ValidateStruct(company); err != nil {
		return errors.BadRequest.New(err.Error())
	}

	return nil
}

func (p *Company) deleteCompany(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "please send a company id"})
	}

	err := p.deleteCompanyUseCase.Execute(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "company deleted"})
}

func (p *Company) updateCompany(c echo.Context) error {
	var company domain.Company
	if err := c.Bind(&company); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "please send a company object"})
	}

	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "please send a company id"})
	}

	if err := company.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	err := p.validateCompany(&company)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	err = p.updateCompanyUseCase.Execute(c.Request().Context(), id, &company)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "company updated"})
}
