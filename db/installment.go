package db

import "aspire/models"

func (db *DB) FindInstallmentById(id string) (*models.Installment, error) {

	return &models.Installment{}, nil
}

func (db *DB) FindInstallments(username *string) ([]*models.Installment, error) {
	return []*models.Installment{}, nil

}

func (db *DB) AddInstallments(installments []*models.Installment, loanId int64) ([]*models.Installment, error) {
	responseInstallment := []*models.Installment{}
	for _, installment := range installments {
		installment.LoanID = loanId

		/*			loan.Currency = "USD"
				}
				loan.State = "PENDING"
				sqlInsert := `
		INSERT INTO loan(amount, term, currency,state)
		VALUES ($1, $2, $3, $4)
		RETURNING id`
				var id int64
				err := db.QueryRow(sqlInsert, loan.Amount, loan.Term, loan.Currency, loan.State).Scan(&id)
				if err != nil {
					return responseLoan, err
				}
				loan.ID = &id
				responseLoan = append(responseLoan, loan)*/
	}
	return responseInstallment, nil
}

func (db *DB) AddInstallment(loan *models.Loan) (*models.Loan, error) {
	return &models.Loan{}, nil
}

func (db *DB) DeleteInstallment(id string) (*models.Loan, error) {
	return &models.Loan{}, nil
}
