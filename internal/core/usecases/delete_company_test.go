package usecases

import (
	"context"
	"errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	mock_kafka "github.com/rhuandantas/xm-challenge/test/mock/kafka"
	mock_mysql "github.com/rhuandantas/xm-challenge/test/mock/repository"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Delete company unit tests", func() {
	var (
		ctx             = context.Background()
		mockCtrl        *gomock.Controller
		repo            *mock_mysql.MockCompanyRepo
		producer        *mock_kafka.MockProducer
		deleteCompanyUC DeleteCompany
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		repo = mock_mysql.NewMockCompanyRepo(mockCtrl)
		producer = mock_kafka.NewMockProducer(mockCtrl)
		deleteCompanyUC = NewDeleteCompany(repo, producer)
	})

	Context("deleting a company", func() {
		It("should return nil", func() {
			repo.EXPECT().DeleteByID(ctx, "company").Return(nil)
			producer.EXPECT().Produce(ctx, "company-events", "message", gomock.Any()).Return(nil)
			err := deleteCompanyUC.Execute(ctx, "company")
			Expect(err).To(BeNil())
		})
		When("get an error from deleting company by name", func() {
			It("should return error", func() {
				repo.EXPECT().DeleteByID(ctx, "company").Return(errors.New("error deleting company by name"))
				err := deleteCompanyUC.Execute(ctx, "company")
				Expect(err).To(HaveOccurred())
			})
		})
		When("get an error from producing message", func() {
			It("should return nil", func() {
				repo.EXPECT().DeleteByID(ctx, "company").Return(nil)
				producer.EXPECT().Produce(ctx, "company-events", "message", gomock.Any()).Return(errors.New("error producing message"))
				err := deleteCompanyUC.Execute(ctx, "company")
				Expect(err).To(BeNil())
			})
		})
	})
})
