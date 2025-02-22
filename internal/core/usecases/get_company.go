package usecases

import (
	"context"
	"github.com/rhuandantas/xm-challenge/internal/adapters/repository"
	"github.com/rhuandantas/xm-challenge/internal/core/domain"
	"github.com/rs/zerolog/log"
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
	log.Info().Msgf("getting company by name %s", name)
	return s.repo.GetByName(ctx, name)
}
