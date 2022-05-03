package db

import (
	"aspire/constants"
	"aspire/models"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-openapi/strfmt"
)

func (db *DB) FindInstallmentById(id string) (*models.Installment, error) {
	installment := &models.Installment{}
	sqlFindById := `
		SELECT id,installment_amount,due_date,state,loan_id,repayment_date,created_on,modified_on from installments where id=$1`
	err := db.QueryRow(sqlFindById, id).Scan(installment.ID, installment.InstallmentAmount, installment.DueDate, installment.State, installment.LoanID, installment.RepaymentTime, installment.CreatedOn, installment.ModifiedOn)
	if err != nil {
		return nil, err
	}
	return installment, nil
}

func (db *DB) FindInstallments(loanId int64, state string, sort string, limit int64, page int64) ([]*models.Installment, error) {
	installments := []*models.Installment{}

	if sort == "" {
		sort = "createdOn.desc"
	}
	var columnKey, operatorKey, columnName, operator string

	sortKeys := strings.Split(sort, ".")
	if len(sortKeys) == 2 {
		columnKey = sortKeys[0]
		operatorKey = sortKeys[1]
	} else {
		return nil, errors.New("sort value not proper")
	}

	switch columnKey {
	case "createdOn":
		columnName = "created_on"
		break
	case "installmentamount":
		columnName = "installment_amount"
		break
	case "dueDate":
		columnName = "due_date"
		break
	case "state":
		columnName = "state"
	}

	switch operatorKey {
	case "asc":
		operator = "ASC"
	case "desc":
		operator = "DESC"
	}

	sortClause := fmt.Sprintf("ORDER BY %s %s", columnName, operator)

	sqlFind := `
		SELECT id,installment_amount,due_date,state,loan_id,repayment_date,created_on,modified_on from installments`
	sqlFindByLoan := `
		SELECT id,installment_amount,due_date,state,loan_id,repayment_date,created_on,modified_on from installments where user_id=$1`
	sqlFindByState := `
		SELECT id,installment_amount,due_date,state,loan_id,repayment_date,created_on,modified_on from installments where state=$1`
	sqlFindByLoanAndState := `
	SELECT id,installment_amount,due_date,state,loan_id,repayment_date,created_on,modified_on from installments where user_id=$1 and state=$2`
	var rows *sql.Rows
	var err error

	if loanId != 0 {
		if state != "" {
			sqlFindByLoanAndState = sqlFindByLoanAndState + " " + sortClause
			if limit != 0 && page != 0 {
				sqlFindByLoanAndState = sqlFindByLoanAndState + " " + "LIMIT $3 " + "OFFSET $4"
				rows, err = db.Query(sqlFindByLoanAndState, loanId, state, limit, page*limit)
				if err != nil {
					return nil, err
				}
			} else {
				rows, err = db.Query(sqlFindByLoanAndState, loanId, state)
				if err != nil {
					return nil, err
				}
			}

		} else {
			sqlFindByLoan = sqlFindByLoan + " " + sortClause
			if limit != 0 && page != 0 {
				sqlFindByLoan = sqlFindByLoan + " " + "LIMIT $2 " + "OFFSET $3"
				rows, err = db.Query(sqlFindByLoan, loanId, limit, page*limit)
				if err != nil {
					return nil, err
				}
			} else {

				rows, err = db.Query(sqlFindByLoan, loanId)
				if err != nil {
					return nil, err
				}
			}
		}
	} else {
		if state != "" {
			sqlFindByState = sqlFindByState + " " + sortClause
			if limit != 0 && page != 0 {
				sqlFindByState = sqlFindByState + " " + "LIMIT $2 " + "OFFSET $3"
				rows, err = db.Query(sqlFindByState, state, limit, page*limit)
				if err != nil {
					return nil, err
				}
			} else {
				rows, err = db.Query(sqlFindByState, state)
				if err != nil {
					return nil, err
				}
			}

		} else {
			sqlFind = sqlFind + " " + sortClause
			if limit != 0 && page != 0 {
				sqlFind = sqlFind + " " + "LIMIT $1 " + "OFFFSET $2"
				rows, err = db.Query(sqlFind, limit, page*limit)
				if err != nil {
					return nil, err
				}
			} else {
				rows, err = db.Query(sqlFind)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	defer rows.Close()
	for rows.Next() {
		installment := &models.Installment{}
		if err := rows.Scan(installment.ID, installment.InstallmentAmount, installment.DueDate, installment.State, installment.LoanID, installment.RepaymentTime, installment.CreatedOn, installment.ModifiedOn); err != nil {
			return nil, err
		}
		installments = append(installments, installment)
	}
	return installments, nil
}

func (db *DB) AddInstallment(installment *models.Installment, loanId int64) (*models.Installment, error) {
	currentTime := time.Now().UTC()
	defaultTime := time.Time{}.UTC()
	installment.LoanID = loanId
	installment.CreatedOn = strfmt.DateTime(currentTime)
	installment.ModifiedOn = strfmt.DateTime(currentTime)
	installment.State = constants.InstallmentStatusPending
	installment.RepaymentAmount = 0
	installment.RepaymentTime = strfmt.DateTime(defaultTime)
	sqlInsert := `
	INSERT INTO installments(installment_amount,due_date,state,loan_id,repayment_time,created_on,modified_on)
	VALUES ($1, $2, $3, $4,$5)
	RETURNING id`
	var id int64
	err := db.QueryRow(sqlInsert, installment.InstallmentAmount, installment.DueDate, installment.State, installment.LoanID, installment.RepaymentTime, installment.CreatedOn, installment.ModifiedOn).Scan(&id)
	if err != nil {
		return nil, err
	}
	installment.ID = &id
	return installment, nil
}

func (db *DB) DeleteInstallment(id int64) error {
	sqlDelete := `
	DELETE FROM installments WHERE id = $1;`
	_, err := db.Exec(sqlDelete, id)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) RepayInstallment(installmentId int64, repaymentAmount float64) error {
	sqlUpdate := `
		UPDATE installments
		SET state = $1, repayment_amount =$2,repayment_time = $3, modifiedOn = $4
		WHERE id = $5;`

	currentTime := time.Now().UTC()
	_, err := db.Exec(sqlUpdate, constants.InstallmentStatusPaid, repaymentAmount, strfmt.DateTime(currentTime), strfmt.DateTime(currentTime), installmentId)
	if err != nil {
		return err
	}
	return nil
}
