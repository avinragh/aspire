package db

import (
	"aspire/models"
)

func (db *DB) FindLoanById(id string) (*models.Loan, error) {
	return &models.Loan{}, nil
}

func (db *DB) FindLoans(username *string) ([]*models.Loan, error) {
	return []*models.Loan{}, nil

}

func (db *DB) AddLoans(loans []*models.Loan) ([]*models.Loan, error) {
	return []*models.Loan{}, nil
}

func (db *DB) AddLoan(loan *models.Loan) (*models.Loan, error) {
	return &models.Loan{}, nil
}

func (db *DB) DeleteLoan(id string) (*models.Loan, error) {
	return &models.Loan{}, nil
}
