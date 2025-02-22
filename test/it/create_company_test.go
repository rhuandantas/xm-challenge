package it

import (
	"context"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rhuandantas/xm-challenge/config"
	"github.com/rhuandantas/xm-challenge/internal/adapters/messaging/kafka"
	"github.com/rhuandantas/xm-challenge/internal/adapters/repository"
	"github.com/rhuandantas/xm-challenge/internal/adapters/repository/mysql"
	"github.com/rhuandantas/xm-challenge/internal/core/domain"
	"github.com/rhuandantas/xm-challenge/internal/core/usecases"
	"testing"
)

func Test(t *testing.T) {
	t.Parallel()
	RegisterFailHandler(Fail)
	RunSpecs(t, "suite test")
}

var _ = Describe("Company Integration Test", func() {
	var (
		ctx                  context.Context
		mysqlContainer       *TestContainer
		kafkaContainer       *TestContainer
		createCompanyUseCase usecases.CreateCompany
		getCompanyUseCase    usecases.GetCompany
		updateCompanyUseCase usecases.UpdateCompany
		deleteCompanyUseCase usecases.DeleteCompany
	)

	BeforeEach(func() {
		ctx = context.Background()
		var err error

		mysqlContainer, err = SetupMySQLContainer(ctx)
		Expect(err).NotTo(HaveOccurred())

		kafkaContainer, err = SetupKafkaContainer(ctx)
		Expect(err).NotTo(HaveOccurred())

		createCompanyUseCase = newCreateCompanyUC(mysqlContainer, kafkaContainer)
		getCompanyUseCase = newGetCompanyUC(mysqlContainer)
		updateCompanyUseCase = newUpdateCompanyUC(mysqlContainer, kafkaContainer)
		deleteCompanyUseCase = newDeleteCompanyUC(mysqlContainer, kafkaContainer)
	})

	AfterEach(func() {
		if mysqlContainer != nil {
			Expect(mysqlContainer.Container.Terminate(ctx)).To(Succeed())
		}
		if kafkaContainer != nil {
			Expect(kafkaContainer.Container.Terminate(ctx)).To(Succeed())
		}
	})

	description := "This is a test company"
	company := &domain.Company{
		Name:              "Test Company",
		Description:       &description,
		AmountOfEmployees: 100,
		Registered:        true,
		Type:              domain.Corporations,
	}

	It("end to end", func() {
		result, err := createCompanyUseCase.Execute(ctx, company)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(BeNil())

		By("try to create the same company again")
		_, err = createCompanyUseCase.Execute(ctx, company)
		Expect(err).To(HaveOccurred()) // should return an error
		By("get the company from the database")
		company, err := getCompanyUseCase.Execute(ctx, company.Name)
		Expect(err).NotTo(HaveOccurred())
		Expect(company).NotTo(BeNil())
		Expect(company.Name).To(Equal("Test Company"))
		Expect(company.Description).To(Equal(&description))
		Expect(company.AmountOfEmployees).To(Equal(100))
		Expect(company.Registered).To(BeTrue())
		Expect(company.Type).To(Equal(domain.Corporations))

		By("update the company")
		err = updateCompanyUseCase.Execute(ctx, company.Id, &domain.Company{
			Name:              "New Company",
			Description:       &description,
			AmountOfEmployees: 200,
			Registered:        false,
			Type:              domain.NonProfit,
		})
		Expect(err).NotTo(HaveOccurred())

		By("get the updated company from the database")
		company, err = getCompanyUseCase.Execute(ctx, "New Company")
		Expect(err).NotTo(HaveOccurred())
		Expect(company).NotTo(BeNil())
		Expect(company.Name).To(Equal("New Company"))
		Expect(company.Description).To(Equal(&description))
		Expect(company.AmountOfEmployees).To(Equal(200))
		Expect(company.Type).To(Equal(domain.NonProfit))

		By("try to delete the company")
		err = deleteCompanyUseCase.Execute(ctx, company.Id)
		Expect(err).NotTo(HaveOccurred())

		By("try to get the company again")
		company, err = getCompanyUseCase.Execute(ctx, "New Company")
		Expect(err).To(HaveOccurred())
		Expect(company).To(BeNil())
	})
})

func newDeleteCompanyUC(musqlC *TestContainer, kafkaC *TestContainer) usecases.DeleteCompany {
	cfg := &config.Config{
		Kafka: config.KafkaConfig{
			Brokers: kafkaC.Addresses,
		},
		Database: config.DatabaseConfig{
			Url: fmt.Sprintf("root:root@tcp(%s:%s)/test_db", musqlC.Host, musqlC.Port),
		},
	}
	producer, err := kafka.NewProducer(cfg)
	if err != nil {
		panic(err)
	}

	dbConn := mysql.NewMySQLConnector(cfg)
	repo := repository.NewCompanyRepo(dbConn)

	return usecases.NewDeleteCompany(repo, producer)
}

func newUpdateCompanyUC(mysqlC *TestContainer, kafkaC *TestContainer) usecases.UpdateCompany {
	cfg := &config.Config{
		Kafka: config.KafkaConfig{
			Brokers: kafkaC.Addresses,
		},
		Database: config.DatabaseConfig{
			Url: fmt.Sprintf("root:root@tcp(%s:%s)/test_db", mysqlC.Host, mysqlC.Port),
		},
	}
	producer, err := kafka.NewProducer(cfg)
	if err != nil {
		panic(err)
	}

	dbConn := mysql.NewMySQLConnector(cfg)
	repo := repository.NewCompanyRepo(dbConn)

	return usecases.NewUpdateCompany(repo, producer)
}

func newGetCompanyUC(mysqlC *TestContainer) usecases.GetCompany {
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Url: fmt.Sprintf("root:root@tcp(%s:%s)/test_db", mysqlC.Host, mysqlC.Port),
		},
	}
	dbConn := mysql.NewMySQLConnector(cfg)
	repo := repository.NewCompanyRepo(dbConn)

	return usecases.NewGetCompany(repo)
}

func newCreateCompanyUC(mysqlC, kafkaC *TestContainer) usecases.CreateCompany {
	cfg := &config.Config{
		Kafka: config.KafkaConfig{
			Brokers: kafkaC.Addresses,
		},
		Database: config.DatabaseConfig{
			Url: fmt.Sprintf("root:root@tcp(%s:%s)/test_db", mysqlC.Host, mysqlC.Port),
		},
	}
	producer, err := kafka.NewProducer(cfg)
	if err != nil {
		panic(err)
	}

	dbConn := mysql.NewMySQLConnector(cfg)
	repo := repository.NewCompanyRepo(dbConn)

	return usecases.NewCreateCompany(repo, producer)
}
