package usecases

import (
	"context"
	"github.com/labstack/gommon/log"
	"github.com/rhuandantas/xm-challenge/internal/adapters/messaging/kafka"
	"github.com/rhuandantas/xm-challenge/internal/adapters/repository"
)

type DeleteCompany interface {
	Execute(ctx context.Context, id string) error
}

type deleteCompany struct {
	repo     repository.CompanyRepo
	producer kafka.Producer
}

func NewDeleteCompany(repo repository.CompanyRepo, producer kafka.Producer) DeleteCompany {
	return &deleteCompany{
		repo:     repo,
		producer: producer,
	}
}

func (s *deleteCompany) Execute(ctx context.Context, id string) error {
	if err := s.repo.DeleteByID(ctx, id); err != nil {
		return err
	}

	if err := s.producer.Produce(ctx, "company-events", "message", map[string]interface{}{"action": "delete", "payload": id}); err != nil {
		log.Warn("error producing message ", err.Error())
	}

	return nil
}
