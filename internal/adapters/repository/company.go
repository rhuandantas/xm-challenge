package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/joomcode/errorx"
	"github.com/rhuandantas/xm-challenge/internal/adapters/repository/mysql"
	"github.com/rhuandantas/xm-challenge/internal/core/domain"
	_error "github.com/rhuandantas/xm-challenge/internal/errors"
	"xorm.io/xorm"
)

const tableName string = "company"

var (
	ErrCompanyNotFound = _error.NotFound.New("company not found")
)

type company struct {
	session *xorm.Session
}

type CompanyRepo interface {
	Create(ctx context.Context, company *domain.Company) (*domain.Company, error)
	GetByName(ctx context.Context, name string) (*domain.Company, error)
	GetByID(ctx context.Context, id string) (*domain.Company, error)
	Update(ctx context.Context, id string, company *domain.Company) (*domain.Company, error)
	DeleteByID(ctx context.Context, id string) error
}

func NewCompanyRepo(connector mysql.DBConnector) CompanyRepo {
	err := connector.SyncTables(new(domain.Company))
	if err != nil {
		panic(err)
	}
	session := connector.GetORM().Table(tableName)
	return &company{
		session: session,
	}
}

func (r *company) Create(ctx context.Context, company *domain.Company) (*domain.Company, error) {
	pkUUID := uuid.New()
	company.Id = pkUUID.String()

	_, err := r.session.Context(ctx).Insert(company)
	if err != nil {
		if _error.IsDuplicatedEntryError(err) {
			return nil, _error.BadRequest.New("company already exists")
		}
		return nil, err
	}

	return company, nil
}

func (r *company) GetByName(ctx context.Context, name string) (*domain.Company, error) {
	companyFilter := &domain.Company{
		Name: name,
	}
	found, err := r.session.Context(ctx).Get(companyFilter)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, ErrCompanyNotFound
	}

	return companyFilter, nil
}

func (r *company) GetByID(ctx context.Context, id string) (*domain.Company, error) {
	companyFilter := &domain.Company{
		Id: id,
	}
	found, err := r.session.Context(ctx).Get(companyFilter)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, ErrCompanyNotFound
	}

	return companyFilter, nil
}

func (r *company) Update(ctx context.Context, id string, company *domain.Company) (*domain.Company, error) {
	_, err := r.session.Context(ctx).ID(id).Update(company)
	if err != nil {
		return nil, err
	}

	return company, nil
}

func (r *company) DeleteByID(ctx context.Context, id string) error {
	_, err := r.session.Context(ctx).Delete(&domain.Company{Id: id})
	if err != nil {
		return errorx.InternalError.New(err.Error())
	}
	return nil
}
