package usecases

import (
	"context"
	"errors"
	"github.com/labstack/gommon/log"
	"github.com/rhuandantas/xm-challenge/internal/adapters/async/kafka"
	"github.com/rhuandantas/xm-challenge/internal/adapters/repository"
	"github.com/rhuandantas/xm-challenge/internal/core/domain"
)

type UpdateCompany interface {
	Execute(ctx context.Context, id string, company *domain.Company) error
}

type updateCompany struct {
	repo     repository.CompanyRepo
	producer kafka.Producer
}

func NewUpdateCompany(repo repository.CompanyRepo, producer kafka.Producer) UpdateCompany {
	return &updateCompany{
		repo:     repo,
		producer: producer,
	}
}

func (p *updateCompany) Execute(ctx context.Context, id string, company *domain.Company) error {
	found, err := p.repo.GetByID(ctx, id)
	if err != nil {
		if !errors.Is(err, repository.ErrCompanyNotFound) {
			return err
		}
	}

	if found == nil {
		return errors.New("company not found")
	}

	company.Id = id
	_, err = p.repo.Update(ctx, id, company)
	if err != nil {
		return errors.New("error creating company")
	}

	if err := p.producer.Produce(ctx, "company-events", "message", map[string]interface{}{"action": "update", "payload": company}); err != nil {
		log.Warn("error producing message ", err.Error())
	}
	return nil
}
