package usecases

import (
	"context"
	"errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	mock_mysql "github.com/rhuandantas/xm-challenge/test/mock/repository"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Get company unit tests", func() {
	var (
		ctx          = context.Background()
		mockCtrl     *gomock.Controller
		repo         *mock_mysql.MockCompanyRepo
		getCompanyUC GetCompany
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		repo = mock_mysql.NewMockCompanyRepo(mockCtrl)
		getCompanyUC = NewGetCompany(repo)
	})

	Context("When getting a company", func() {
		It("should return nil", func() {
			repo.EXPECT().GetByName(ctx, "company").Return(nil, nil)
			_, err := getCompanyUC.Execute(ctx, "company")
			Expect(err).To(BeNil())
		})
		When("get an error from getting company by name", func() {
			It("should return error", func() {
				repo.EXPECT().GetByName(ctx, "company").Return(nil, errors.New("error getting company by name"))
				_, err := getCompanyUC.Execute(ctx, "company")
				Expect(err).To(HaveOccurred())
			})
		})
	})

})
