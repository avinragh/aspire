package db

import (
	"aspire/constants"
	"aspire/models"
	"aspire/util"
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/go-openapi/strfmt"
)

func (db *DB) FindLoanById(id int64) (*models.Loan, error) {
	loan := &models.Loan{}
	sqlFindById := `
		SELECT id,amount,created_on,term,currency,state,modified_on,start_date,user_Id FROM loans WHERE id=$1`
	err := db.QueryRow(sqlFindById, id).Scan(loan.ID, loan.Amount, loan.CreatedOn, loan.Term, loan.Currency, loan.State, loan.ModifiedOn, loan.StartDate, loan.UserID)
	if err != nil {
		return nil, err
	}
	return loan, nil
}

func (db *DB) FindLoans(userId *int64) ([]*models.Loan, error) {
	loans := []*models.Loan{}
	sqlFind := `
		SELECT id,amount,created_on,term,currency,state,modified_on,start_date,user_Id FROM loans`
	sqlFindByUser := `
		SELECT id,amount,created_on,term,currency,state,modified_on,start_date,user_Id FROM loans where user_id=$1`
	var rows *sql.Rows
	var err error
	if userId != nil {
		rows, err = db.Query(sqlFindByUser, *userId)
		if err != nil {
			return nil, err
		}
	} else {
		rows, err = db.Query(sqlFind)
		if err != nil {
			return nil, err
		}

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
	defaultTime := strfmt.DateTime(time.Time{}.UTC())
	currentDate := strfmt.DateTime(time.Now().UTC())
	loan.CreatedOn = currentDate
	loan.ModifiedOn = currentDate
	loan.UserID = userId
	sqlInsert := `
		INSERT INTO loans(amount,created_on,term,currency,state,modified_on,start_date,user_id)
		VALUES ($1, $2, $3, $4,$5, $6, $7, $8)
		RETURNING id`
	var id int64
	err := db.QueryRow(sqlInsert, loan.Amount, loan.CreatedOn, loan.Term, loan.Currency, loan.State, loan.ModifiedOn, defaultTime, loan.UserID).Scan(&id)
	if err != nil {
		return nil, err
	}
	loan.ID = &id
	return loan, nil
}

func (db *DB) ApproveLoan(loanId int64, installments []*models.Installment) error {
	sqlUpdate := `
		UPDATE loans
		SET state = $1, modifiedOn = $2, startDate = $3
		WHERE id= $4;`
	sqlInsert := `
		INSERT INTO installments(installment_amount,due_date,state,loan_id,created_on,modified_on)
		VALUES ($1, $2, $3, $4, $5)`

	currentTime := time.Now().UTC()
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	_, err = tx.ExecContext(ctx, sqlUpdate, constants.LoanStatusApproved, strfmt.DateTime(currentTime), strfmt.DateTime(util.GetDate(currentTime)), loanId)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, installment := range installments {
		installment.CreatedOn = strfmt.DateTime(currentTime)
		installment.ModifiedOn = strfmt.DateTime(currentTime)
		_, err = tx.ExecContext(ctx, sqlInsert, installment.InstallmentAmount, installment.DueDate, installment.State, loanId, installment.CreatedOn, installment.ModifiedOn)
		if err != nil {
			tx.Rollback()
			return err
		}

	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) DeleteLoan(id int64) error {
	sqlDelete := `
	DELETE FROM loans WHERE id=$1;`
	_, err := db.Exec(sqlDelete, id)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) UpdateToPaid(id int64) error {
	currentDate := time.Now().UTC()
	sqlUpdate := `
	UPDATE loans
	SET state= $1,modifiedOn =$2`
	_, err := db.Exec(sqlUpdate, constants.LoanStatusPaid, strfmt.DateTime(currentDate))
	if err != nil {
		return err
	}
	return nil
}
