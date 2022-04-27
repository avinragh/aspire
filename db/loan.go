package db

import (
	"aspire/models"
	"time"

	"github.com/go-openapi/strfmt"
)

func (db *DB) FindLoanById(id int64) (*models.Loan, error) {
	loan := &models.Loan{}
	sqlFindById := `
		SELECT id,amount,createdOn,term,currency,state,modifiedOn,startDate,userId FROM loans WHERE id=$1`
	err := db.QueryRow(sqlFindById, id).Scan(loan.ID, loan.Amount, loan.CreatedOn, loan.Term, loan.Currency, loan.State, loan.ModifiedOn, loan.StartDate, loan.UserID)
	if err != nil {
		return nil, err
	}
	return loan, nil
}

func (db *DB) FindLoans() ([]*models.Loan, error) {
	loans := []*models.Loan{}
	sqlFind := `
		SELECT id,amount,createdOn,term,currency,state,modifiedOn,startDate,userId FROM loans`
	rows, err := db.Query(sqlFind)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		loan := &models.Loan{}
		if err := rows.Scan(loan.ID, loan.Amount, loan.CreatedOn, loan.Term, loan.Currency, loan.State, loan.ModifiedOn, loan.StartDate, loan.UserID); err != nil {
			return nil, err
		}
		loans = append(loans, loan)
	}

	return loans, nil

}

func (db *DB) AddLoan(loan *models.Loan, userId int64) (*models.Loan, error) {
	defaultTime := strfmt.DateTime(time.Time{})
	currentDate := strfmt.DateTime(time.Now())
	sqlInsert := `
		INSERT INTO loans(amount,createdOn,term,currency,state,modifiedOn,startDate,userId)
		VALUES ($1, $2, $3, $4,$5,$6,$7,$8)
		RETURNING id`
	var id int64
	err := db.QueryRow(sqlInsert, loan.Amount, currentDate, loan.Term, loan.Currency, loan.State, currentDate, defaultTime, userId).Scan(&id)
	if err != nil {
		return nil, err
	}
	loan.ID = &id
	return loan, nil
}

/*func (db *DB) AddLoan(loan *models.Loan) (*models.Loan, error) {
	return &models.Loan{}, nil
}

func (db *DB) UpdateLoan(loan *models.Loan) (*models.Loan, error) {

}

func (db *DB) DeleteLoan(id string) (*models.Loan, error) {
	return &models.Loan{}, nil
}*/
