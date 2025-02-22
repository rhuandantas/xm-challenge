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

var _ = Describe("Create company unit tests", func() {
	var (
		ctx             = context.Background()
		mockCtrl        *gomock.Controller
		repo            *mock_mysql.MockCompanyRepo
		producer        *mock_kafka.MockProducer
		createCompanyUC CreateCompany
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		repo = mock_mysql.NewMockCompanyRepo(mockCtrl)
		producer = mock_kafka.NewMockProducer(mockCtrl)
		createCompanyUC = NewCreateCompany(repo, producer)
	})

	Context("When creating a company", func() {
		It("should return nil", func() {
			repo.EXPECT().GetByName(ctx, "company").Return(nil, nil)
			repo.EXPECT().Create(ctx, gomock.Any()).Return(nil, nil)
			producer.EXPECT().Produce(ctx, "company-events", "message", gomock.Any()).Return(nil)
			_, err := createCompanyUC.Execute(ctx, &domain.Company{Name: "company"})
			Expect(err).To(BeNil())
		})
		When("get an error from getting company by name", func() {
			It("should return error", func() {
				repo.EXPECT().GetByName(ctx, "company").Return(nil, errors.New("error getting company by name"))
				_, err := createCompanyUC.Execute(ctx, &domain.Company{Name: "company"})
				Expect(err).To(HaveOccurred())
			})
		})
		When("company already exists", func() {
			It("should return error", func() {
				repo.EXPECT().GetByName(ctx, "company").Return(&domain.Company{}, nil)
				_, err := createCompanyUC.Execute(ctx, &domain.Company{Name: "company"})
				Expect(err).To(HaveOccurred())
			})
		})
		When("get an error from creating company", func() {
			It("should return error", func() {
				repo.EXPECT().GetByName(ctx, "company").Return(nil, nil)
				repo.EXPECT().Create(ctx, gomock.Any()).Return(nil, errors.New("error creating company"))
				_, err := createCompanyUC.Execute(ctx, &domain.Company{Name: "company"})
				Expect(err).To(HaveOccurred())
			})
		})
		When("get an error from producing message", func() {
			It("should return nil", func() {
				repo.EXPECT().GetByName(ctx, "company").Return(nil, nil)
				repo.EXPECT().Create(ctx, gomock.Any()).Return(nil, nil)
				producer.EXPECT().Produce(ctx, "company-events", "message", gomock.Any()).Return(errors.New("error producing message"))
				_, err := createCompanyUC.Execute(ctx, &domain.Company{Name: "company"})
				Expect(err).To(BeNil())
			})
		})
	})

})
