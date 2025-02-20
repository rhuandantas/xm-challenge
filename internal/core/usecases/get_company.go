package usecases

import (
	"context"
	"github.com/rhuandantas/xm-challenge/internal/adapters/repository"
	"github.com/rhuandantas/xm-challenge/internal/core/domain"
)

type GetCompany interface {
	Execute(ctx context.Context, name string) (*domain.Company, error)
}

type getCompany struct {
	repo repository.CompanyRepo
}

func NewGetCompany(repo repository.CompanyRepo) GetCompany {
	return &getCompany{
		repo: repo,
	}
}

func (s *getCompany) Execute(ctx context.Context, name string) (*domain.Company, error) {
	return s.repo.GetByName(ctx, name)
}
