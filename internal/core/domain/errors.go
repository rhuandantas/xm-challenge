package domain

import "errors"

var (
	ErrInvalidCompanyType    = errors.New("invalid company type. Must be one of: corporations, nonprofit, cooperative, sole proprietorship")
	ErrEmptyCompanyName      = errors.New("company name is invalid")
	ErrAmountOfEmployeesZero = errors.New("amount_of_employee cannot be zero")
)
