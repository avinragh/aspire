package db

import (
	"aspire/constants"
	"aspire/models"
)

func (db *DB) FindLoanById(id string) (*models.Loan, error) {

	return &models.Loan{}, nil
}

func (db *DB) FindLoans(username *string) ([]*models.Loan, error) {
	return []*models.Loan{}, nil

}

func (db *DB) AddLoan(loan *models.Loan) (*models.Loan, error) {
	if loan.Currency == "" {
		loan.Currency = "USD"
	}
	loan.State = constants.LoanStatusPending
	sqlInsert := `
INSERT INTO loan(amount, term, currency,state)
VALUES ($1, $2, $3, $4)
RETURNING id`
	var id int64
	err := db.QueryRow(sqlInsert, loan.Amount, loan.Term, loan.Currency, loan.State).Scan(&id)
	if err != nil {
		return nil, err
	}
	loan.ID = &id
	return loan, nil
}

/*func (db *DB) AddLoan(loan *models.Loan) (*models.Loan, error) {
	return &models.Loan{}, nil
}*/

func (db *DB) DeleteLoan(id string) (*models.Loan, error) {
	return &models.Loan{}, nil
}
