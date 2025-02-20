package usecases

import (
	"context"
	"errors"
	"github.com/labstack/gommon/log"
	"github.com/rhuandantas/xm-challenge/internal/adapters/async/kafka"
	"github.com/rhuandantas/xm-challenge/internal/adapters/repository"
	"github.com/rhuandantas/xm-challenge/internal/core/domain"
)

type CreateCompany interface {
	Execute(ctx context.Context, company *domain.Company) (*domain.Company, error)
}

type createCompany struct {
	repo     repository.CompanyRepo
	producer kafka.Producer
}

func NewCreateCompany(repo repository.CompanyRepo, producer kafka.Producer) CreateCompany {
	return &createCompany{
		repo:     repo,
		producer: producer,
	}
}

func (p *createCompany) Execute(ctx context.Context, company *domain.Company) (*domain.Company, error) {
	found, err := p.repo.GetByName(ctx, company.Name)
	if err != nil {
		if !errors.Is(err, repository.ErrCompanyNotFound) {
			return nil, err
		}
	}

	if found != nil {
		return nil, errors.New("company already exists")
	}

	_, err = p.repo.Create(ctx, company)
	if err != nil {
		return nil, errors.New("error creating company")
	}

	err = p.producer.Produce(ctx, "company-events", "message", map[string]interface{}{"action": "create", "payload": company})
	if err != nil {
		log.Warn("error producing message ", err.Error())
	}

	return nil, nil
}
