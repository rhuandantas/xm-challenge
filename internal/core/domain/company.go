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
	if !c.Type.IsValid() {
		return ErrInvalidCompanyType
	}
	return nil
}

type Company struct {
	Id                string      `xorm:"pk uuid"`
	Name              string      `json:"name" xorm:"unique varchar(15) notnull"`
	Description       *string     `json:"description,omitempty" xorm:"varchar(3000)"`
	AmountOfEmployees int         `json:"amount_of_employees" xorm:"notnull"`
	Registered        bool        `json:"registered" xorm:"notnull"`
	Type              CompanyType `json:"type" xorm:"notnull"`
}
