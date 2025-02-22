package usecases

import (
	"context"
	"errors"
	"github.com/rhuandantas/xm-challenge/internal/adapters/messaging/kafka"
	"github.com/rhuandantas/xm-challenge/internal/adapters/repository"
	"github.com/rhuandantas/xm-challenge/internal/core/domain"
	"github.com/rs/zerolog/log"
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
			log.Error().Msgf("error getting company: %s", err.Error())
			return err
		}
	}

	if found == nil {
		log.Warn().Msgf("company not found")
		return errors.New("company not found")
	}

	company.Id = id
	_, err = p.repo.Update(ctx, id, company)
	if err != nil {
		log.Error().Msgf("error updating company: %s", err.Error())
		return errors.New("error creating company")
	}

	if err := p.producer.Produce(ctx, "company-events", "message", map[string]interface{}{"action": "update", "payload": company}); err != nil {
		log.Error().Msgf("error producing message: %s", err.Error())
	}

	return nil
}
