package db

import (
	"aspire/constants"
	"aspire/models"
	"aspire/util"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-openapi/strfmt"
)

func (db *DB) FindInstallmentById(id int64) (*models.Installment, error) {
	installment := &models.Installment{}
	sqlFindById := `
		SELECT id,installment_amount,installments.repayment_amount,due_date,state,loan_id,repayment_time,created_on,modified_on from installments where id=$1`
	err := db.QueryRow(sqlFindById, id).Scan(&installment.ID, &installment.InstallmentAmount, &installment.RepaymentAmount, &installment.DueDate, &installment.State, &installment.LoanID, &installment.RepaymentTime, &installment.CreatedOn, &installment.ModifiedOn)
	if err != nil {
		return nil, err
	}
	return installment, nil
}

func (db *DB) FindInstallments(loanId *int64, state *string, sort *string, limit *int64, page *int64, loanIds []int64) ([]*models.Installment, error) {
	installments := []*models.Installment{}

	if sort == nil {
		sortString := "createdOn.desc"
		sort = &sortString
	}
	var columnKey, operatorKey, columnName, operator string

	sortKeys := strings.Split(*sort, ".")
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
		SELECT installments.id,installments.installment_amount,installments.repayment_amount,installments.due_date,installments.state,installments.loan_id,installments.repayment_time,installments.created_on,installments.modified_on from installments`
	sqlFindByLoan := `
		SELECT installments.id,installments.installment_amount,installments.repayment_amount,installments.due_date,installments.state,installments.loan_id,installments.repayment_time,installments.created_on,installments.modified_on from installments where installments.loan_id=$1`
	sqlFindByState := `
		SELECT installments.id,installments.installment_amount,installments.repayment_amount,installments.due_date,installments.state,installments.loan_id,installments.repayment_time,installments.created_on,installments.modified_on from installments where installments.state=$1`
	sqlFindByLoanAndState := `
		SELECT installments.id,installments.installment_amount,installments.repayment_amount,installments.due_date,installments.state,installments.loan_id,installments.repayment_time,installments.created_on,installments.modified_on from installments where installments.loan_id=$1 and installments.state=$2`
	var rows *sql.Rows
	var err error
	loanIdsClause := util.GetLoanIdsClause(loanIds)

	if loanId != nil {
		if state != nil {
			if loanIdsClause != "" {
				sqlFindByLoanAndState = sqlFindByLoanAndState + " AND installments.loan_id IN " + loanIdsClause
			}

			sqlFindByLoanAndState = sqlFindByLoanAndState + " " + sortClause
			if limit != nil && page != nil {
				sqlFindByLoanAndState = sqlFindByLoanAndState + " " + "LIMIT $3 " + "OFFSET $4"
				rows, err = db.Query(sqlFindByLoanAndState, loanId, state, limit, ((*page)-1)*(*limit))
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
			if loanIdsClause != "" {
				sqlFindByLoan = sqlFindByLoan + " AND installments.loan_id IN " + loanIdsClause
			}

			sqlFindByLoan = sqlFindByLoan + " " + sortClause
			if limit != nil && page != nil {
				sqlFindByLoan = sqlFindByLoan + " " + "LIMIT $2 " + "OFFSET $3"
				rows, err = db.Query(sqlFindByLoan, loanId, limit, ((*page)-1)*(*limit))
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
		if state != nil {
			if loanIdsClause != "" {
				sqlFindByState = sqlFindByState + " AND installments.loan_id IN " + loanIdsClause
			}

			sqlFindByState = sqlFindByState + " " + sortClause
			if limit != nil && page != nil {
				sqlFindByState = sqlFindByState + " " + "LIMIT $2 " + "OFFSET $3"
				rows, err = db.Query(sqlFindByState, state, limit, ((*page)-1)*(*limit))
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

			if loanIdsClause != "" {
				sqlFind = sqlFind + " WHERE installments.loan_id IN " + loanIdsClause
			}

			sqlFind = sqlFind + " " + sortClause
			if limit != nil && page != nil {
				sqlFind = sqlFind + " " + "LIMIT $1 " + "OFFSET $2"
				rows, err = db.Query(sqlFind, limit, ((*page)-1)*(*limit))
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
		if err := rows.Scan(&installment.ID, &installment.InstallmentAmount, &installment.RepaymentAmount, &installment.DueDate, &installment.State, &installment.LoanID, &installment.RepaymentTime, &installment.CreatedOn, &installment.ModifiedOn); err != nil {
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
	installment.RepaymentAmount = util.GetFloat64Pointer(0)
	installment.RepaymentTime = strfmt.DateTime(defaultTime)
	sqlInsert := `
	INSERT INTO installments(installment_amount, repayment_amount,due_date,state,loan_id,repayment_time,created_on,modified_on)
	VALUES ($1, $2, $3, $4,$5)
	RETURNING installments.id,installments.installment_amount,installments.repayment_amount,installments.due_date,installments.state,installments.loan_id,installments.repayment_time,installments.created_on,installments.modified_on`
	err := db.QueryRow(sqlInsert, installment.InstallmentAmount, installment.RepaymentAmount, installment.DueDate, installment.State, installment.LoanID, installment.RepaymentTime, installment.CreatedOn, installment.ModifiedOn).Scan(&installment.ID, &installment.InstallmentAmount, &installment.RepaymentAmount, &installment.DueDate, &installment.State, &installment.LoanID, &installment.RepaymentTime, &installment.CreatedOn, &installment.ModifiedOn)
	if err != nil {
		return nil, err
	}
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

func (db *DB) RepayInstallment(installmentId int64, repaymentAmount float64, loanIds []int64) (*models.Installment, error) {
	loanIdsClause := util.GetLoanIdsClause(loanIds)
	installment := &models.Installment{}
	sqlUpdate := `
		UPDATE installments
		SET state = $1, repayment_amount =$2,repayment_time = $3, modified_on = $4
		WHERE installments.id = $5`
	if loanIdsClause != "" {
		sqlUpdate = sqlUpdate + " AND installments.loan_id 	IN " + loanIdsClause
	}
	sqlUpdate = sqlUpdate +
		` RETURNING installments.id,installments.installment_amount,installments.repayment_amount,installments.due_date,installments.state,installments.loan_id,installments.repayment_time,installments.created_on,installments.modified_on`

	currentTime := time.Now().UTC()
	err := db.QueryRow(sqlUpdate, constants.InstallmentStatusPaid, repaymentAmount, strfmt.DateTime(currentTime), strfmt.DateTime(currentTime), installmentId).Scan(&installment.ID, &installment.InstallmentAmount, &installment.RepaymentAmount, &installment.DueDate, &installment.State, &installment.LoanID, &installment.RepaymentTime, &installment.CreatedOn, &installment.ModifiedOn)
	if err != nil {
		return nil, err
	}
	return installment, nil
}
