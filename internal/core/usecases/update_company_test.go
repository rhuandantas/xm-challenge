package usecases

import (
	"context"
	"errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rhuandantas/xm-challenge/internal/core/domain"
	mock_kafka "github.com/rhuandantas/xm-challenge/test/mock/kafka"
	mock_mysql "github.com/rhuandantas/xm-challenge/test/mock/repository"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Update company unit tests", func() {
	var (
		ctx             = context.Background()
		mockCtrl        *gomock.Controller
		repo            *mock_mysql.MockCompanyRepo
		producer        *mock_kafka.MockProducer
		updateCompanyUC UpdateCompany
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		repo = mock_mysql.NewMockCompanyRepo(mockCtrl)
		producer = mock_kafka.NewMockProducer(mockCtrl)
		updateCompanyUC = NewUpdateCompany(repo, producer)
	})

	Context("updating a company", func() {
		It("should return nil", func() {
			repo.EXPECT().GetByID(ctx, "company").Return(&domain.Company{}, nil)
			repo.EXPECT().Update(ctx, "company", gomock.Any()).Return(nil, nil)
			producer.EXPECT().Produce(ctx, "company-events", "message", gomock.Any()).Return(nil)
			err := updateCompanyUC.Execute(ctx, "company", &domain.Company{})
			Expect(err).To(BeNil())
		})
		When("get an error from getting company by id", func() {
			It("should return error", func() {
				repo.EXPECT().GetByID(ctx, "company").Return(nil, errors.New("error getting company by id"))
				err := updateCompanyUC.Execute(ctx, "company", &domain.Company{})
				Expect(err).To(HaveOccurred())
			})
		})
		When("company not found", func() {
			It("should return error", func() {
				repo.EXPECT().GetByID(ctx, "company").Return(nil, nil)
				err := updateCompanyUC.Execute(ctx, "company", &domain.Company{})
				Expect(err).To(HaveOccurred())
			})
		})
		When("get an error from updating company", func() {
			It("should return error", func() {
				repo.EXPECT().GetByID(ctx, "company").Return(&domain.Company{}, nil)
				repo.EXPECT().Update(ctx, "company", gomock.Any()).Return(nil, errors.New("error updating company"))
				err := updateCompanyUC.Execute(ctx, "company", &domain.Company{})
				Expect(err).To(HaveOccurred())
			})
		})
		When("get an error from producing message", func() {
			It("should return nil", func() {
				repo.EXPECT().GetByID(ctx, "company").Return(&domain.Company{}, nil)
				repo.EXPECT().Update(ctx, "company", gomock.Any()).Return(nil, nil)
				producer.EXPECT().Produce(ctx, "company-events", "message", gomock.Any()).Return(errors.New("error producing message"))
				err := updateCompanyUC.Execute(ctx, "company", &domain.Company{})
				Expect(err).To(BeNil())
			})
		})

	})

})
