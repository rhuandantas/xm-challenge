package domain

type CompanyType string

const (
	Corporations       CompanyType = "corporations"
	NonProfit          CompanyType = "nonprofit"
	Cooperative        CompanyType = "cooperative"
	SoleProprietorship CompanyType = "sole proprietorship"
)

func (ct CompanyType) IsValid() bool {
	switch ct {
	case Corporations, NonProfit, Cooperative, SoleProprietorship:
		return true
	default:
		return false
	}
}

func (c *Company) Validate() error {
	if c.Name == "" && len(c.Name) < 15 {
		return ErrEmptyCompanyName
	}
	if !c.Type.IsValid() {
		return ErrInvalidCompanyType
	}

	if c.AmountOfEmployees <= 0 {
		return ErrAmountOfEmployeesZero
	}
	return nil
}

type Company struct {
	Id                string      `xorm:"pk uuid"`
	Name              string      `json:"name" validate:"required"  xorm:"unique varchar(15) notnull"`
	Description       *string     `json:"description,omitempty" xorm:"varchar(3000)"`
	AmountOfEmployees int         `json:"amount_of_employees" validate:"required" xorm:"notnull"`
	Registered        bool        `json:"registered" validate:"required" xorm:"notnull"`
	Type              CompanyType `json:"type" validate:"required" xorm:"notnull"`
}
